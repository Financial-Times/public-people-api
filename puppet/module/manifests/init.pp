class public_people_api {

  $binary_name = "public-people-api"
  $install_dir = "/usr/local/$binary_name"
  $binary_file = "$install_dir/$binary_name"
  $log_dir = "/var/log/apps"
  $neoURL=hiera("neoURL", "http://localhost:7474/db/data")
  $port=hiera("port", "8080")

  user { $binary_name:
    ensure    => present,
  }

  file {
    $install_dir:
      mode    => "0664",
      ensure  => directory;

    $binary_file:
      ensure  => present,
      source  => "puppet:///modules/$module_name/$binary_name",
      mode    => "0755",
      require => File[$install_dir];

    $log_dir:
      ensure  => directory,
      mode    => "0664"
  }

  service { 'public-people-api':
    ensure => running,
    enable => true,
    binary => $binary_file,
    flags => "blablabl"
  }

#
  # exec { 'restart_app':
  #   command     => "supervisorctl restart $binary_name",
  #   path        => "/usr/bin:/usr/sbin:/bin",
  #   subscribe   => [
  #     File[$binary_file],
  #     Class["${module_name}::supervisord"]
  #   ],
  #   refreshonly => true
  # }
}
