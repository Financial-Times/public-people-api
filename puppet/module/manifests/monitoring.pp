class public_people_api::monitoring {

  $port = "8080"
  $cmd_check_http_json = "/usr/lib64/nagios/plugins/check_http_json.py --host $hostname:$port --path /__health --key_equals \"\$ARG1\$\""
  $nrpe_cmd_check_http_json = '/usr/lib64/nagios/plugins/check_nrpe -H $HOSTNAME$ -c check_http_json -a "$ARG1$"'
  $action_url = 'https://sites.google.com/a/ft.com/technology/TODO/addrunbookurl'
  $config_file = '/etc/nrpe.d/check_http_json.cfg'

  package {
    'argparse':
      ensure  => 'installed',
      provider => pip,
      require  => Package['python-pip'];
  }

  # https://github.com/drewkerrigan/nagios-http-json ; hash: c678dfd518ebc760e42152c2323ccf17e92e5892
  file { '/usr/lib64/nagios/plugins/check_http_json.py':
    ensure          => 'present',
    mode            => 0755,
    source          => "puppet:///modules/$module_name/check_http_json.py",
  }

  file { $config_file:
    ensure          => 'present',
    mode            => 0644,
    content         => "command[check_http_json]=${$cmd_check_http_json}\n"
  }

  exec { 'reload-nrpe-service':
    command         => '/etc/init.d/nrpe reload',
    refreshonly     => true,
    require         => File[$config_file]
  }

  @@nagios_command { "${hostname}_check_http_json":
    command_line => $nrpe_cmd_check_http_json,
    tag => $content_platform_nagios::client::tags_to_apply
  }

  @@nagios_service { "${hostname}_check_http_json_health_1":
    use                 => "generic-service",
    host_name           =>  "${::certname}",
    check_command       => "${hostname}_check_http_json!checks(0).ok,True",
    check_interval      => 1,
    action_url          => $action_url,
    notes_url           => $action_url,
    notes               => "Severity 2 \\n Service unavailable \\n Public People Read GO App's healthchecks are failing. Please check http://${::hostname}:8080/__health \\n\\n",
    service_description => "Check the application healthcheck",
    display_name        => "${hostname}_check_http_json",
    tag                 => $content_platform_nagios::client::tags_to_apply,
  }

  nagios::nrpe_checks::check_tcp {
    "${::certname}/1":
      host          => "localhost",
      port          => 8080,
      notes         => "check ${::certname} [$hostname] listening on HTTP port 8080 ";
  }
}
