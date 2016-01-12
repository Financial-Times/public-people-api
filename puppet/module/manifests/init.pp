class public_people_api {

  $service_name = "public-people-api"
  $user = $service_name
  $group = $user

  $install_dir = "/usr/local/$service_name"
  $binary_file = "$install_dir/$service_name"
  $initd_file = "/etc/init.d/$service_name"
  $start_file = "$install_dir/start.sh"

  $pid_dir = "/var/run/$service_name"
  $pid_file = "${pid_dir}/${service_name}.pid"

  $log_dir = "/var/log/apps/$service_name"
  $startup_log_file = "${log_dir}/${service_name}-startup.log"
  $status_check_url = "http://`hostname`:8080"
  $startup_timeout = 20

  #class { 'common_pp_up': }

  user { $service_name:
    ensure    => present,
  }

  file {
    $install_dir:
      mode    => "0644",
      ensure  => directory,
      owner   => $user,
      group   => $group,
      recurse => remote,
      source  => "puppet:///modules/$module_name";

    $binary_file:
      ensure  => present,
      source  => "puppet:///modules/$module_name/$service_name",
      owner   => $user,
      group   => $group,
      mode    => "0755",
      require => File[$install_dir];

    "/var/log/apps":
      ensure  => directory,
      mode    => "0755";

    $log_dir:
      ensure  => directory,
      owner   => $service_name,
      group   => $service_name,
      mode    => "0755";

    $startup_log_file:
      require    => File [$log_dir],
      ensure     => file,
      owner      => $service_name,
      group      => 'logs',
      mode       => '0640';

    $pid_dir:
      ensure  => directory,
      owner   => $user,
      group   => $group,
      mode    => "0755";

    $start_file:
      mode    => "0755",
      owner   => $user,
      group   => $group,
      content => template("$module_name/start.sh.erb"),
      require => File[$install_dir];

    $initd_file:
      mode    => "0755",
      content => template("$module_name/initd-script.erb");
    }

    service { $service_name:
      require	    => [File [$install_dir,
                           $binary_file,
                           $log_dir,
                           $startup_log_file,
                           $pid_dir,
                           $start_file,
                           $initd_file],
            Exec["add-service-$service_name"]],
      subscribe   => File [$binary_file],
      hasstatus   => true,
      hasrestart  => true,
      ensure	    => running
  }

  exec {'iptables-forward-port-8080-80':
    command => '/sbin/iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-port 8080; service iptables save',
    user    => 'root',
    group   => 'root',
    unless  => "/sbin/iptables -S -t nat | grep -q 'PREROUTING -p tcp -m tcp --dport 80 -j REDIRECT --to-ports 8080' 2>/dev/null"
  }

  exec {"add-service-$service_name":
    command    => "/sbin/chkconfig --add ${$service_name}",
    unless     => "/sbin/chkconfig --list ${$service_name} 2>/dev/null",
    require    => File["$initd_file"];
  }
}
