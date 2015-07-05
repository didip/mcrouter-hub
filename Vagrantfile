# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.define "centos" do |centos|
    centos.vm.box = "metcalfc/centos70-docker"

    centos.vm.provision :shell, path: "installers/centos7.sh"
    centos.vm.network :private_network, type: :static, ip: "192.168.50.241"
  end

  config.ssh.forward_agent = true

  config.vm.synced_folder ENV['GOPATH'], "/go"
end
