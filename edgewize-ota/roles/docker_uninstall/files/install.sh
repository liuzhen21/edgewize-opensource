#!/bin/bash
# Author: Jrohy
# Github: https://github.com/Jrohy/docker-install

offline_file=""



#######color code########
red="31m"      
green="32m"  
yellow="33m" 
blue="36m"
fuchsia="35m"

# cancel centos alias
[[ -f /etc/redhat-release ]] && unalias -a

sysctl_list=(
    "net.ipv4.ip_forward"
    "net.bridge.bridge-nf-call-iptables"
    "net.bridge.bridge-nf-call-ip6tables"
)

color_echo(){
    local color=$1
    echo -e "\033[${color}${@:2}\033[0m"
}


full_path() {
   local pwd=`pwd`
   if [ -d $1 ]; then
      cd $1
   elif [ -f $1 ]; then
      cd `dirname $1`
   else
      cd
   fi
   echo $(cd ..; cd -)
   cd ${pwd} >/dev/null
}

check_file(){
    local file=$1
    if [[ ! -e $file ]];then
        color_echo $red "$file file not exist!\n"
        exit 1
    elif [[ ! -f $file ]];then
        color_echo $red "$file not a file!\n"
        exit 1
    fi

    file_name=$(echo ${file##*/})
    file_path=$(full_path $file)
    if [[ !  $file_name =~ ".tgz" && !  $file_name =~ ".tar.gz" ]];then
        color_echo $red "$file not a tgz file!\n"
        echo -e "please download docker binary file: $(color_echo $fuchsia $download_url)\n"
        exit 1
    fi
}

#######get params#########
while [[ $# > 0 ]];do
    case "$1" in
        -f|--file=)
        offline_file="$2"
        check_file $offline_file
        shift
        ;;
        -s|--standard)
        standard_mode=1
        shift
        ;;
        -h|--help)
        echo "$0 [-h] [-f file]"
        echo "   -f, --file=[file_path]      offline tgz file path"
        echo ""
        echo "Docker binary download link:  $(color_echo $fuchsia $download_url)"
        exit 0
        shift # past argument
        ;; 
        *)
                # unknown option
        ;;
    esac
    shift # past argument or value
done
#############################

check_sys() {
    if [[ -z `command -v systemctl` ]];then
        color_echo ${red} "system must be have systemd!"
        exit 1
    fi
    if [[ -z `uname -m|grep 64` ]];then
        color_echo ${red} "docker only support 64-bit system!"
        exit 1
    fi
    # check os
    if [[ `command -v apt-get` ]];then
        package_manager='apt-get'
    elif [[ `command -v dnf` ]];then
        package_manager='dnf'
    elif [[ `command -v yum` ]];then
        package_manager='yum'
    else
        color_echo $red "Not support OS!"
        exit 1
    fi
}

write_service(){
        mkdir -p /etc/systemd/system/
        cat > /etc/systemd/system/docker.service << EOF
[Unit]
Description=Docker Application Container Engine
Documentation=https://docs.docker.com
After=network-online.target firewalld.service
Wants=network-online.target
 
[Service]
Type=notify
# the default is not to use systemd for cgroups because the delegate issues still
# exists and systemd currently does not support the cgroup feature set required
# for containers run by docker
ExecStart=/usr/bin/dockerd
ExecReload=/bin/kill -s HUP $MAINPID
# Having non-zero Limit*s causes performance problems due to accounting overhead
# in the kernel. We recommend using cgroups to do container-local accounting.
LimitNOFILE=infinity
LimitNPROC=infinity
LimitCORE=infinity
# Uncomment TasksMax if your systemd version supports it.
# Only systemd 226 and above support this version.
#TasksMax=infinity
TimeoutStartSec=0
# set delegate yes so that systemd does not reset the cgroups of docker containers
Delegate=yes
# kill only the docker process, not all processes in the cgroup
KillMode=process
# restart the docker process if it exits prematurely
Restart=on-failure
StartLimitBurst=3
StartLimitInterval=60s
 
[Install]
WantedBy=multi-user.target
EOF
}

offline_install(){
    local origin_path=$(pwd)
    cd $file_path
    tar xzvf $file_name
    cp -rf docker/* /usr/bin/
    rm -rf docker
}

set_sysctl(){
    for conf in ${sysctl_list[@]}
    do
        check=`sysctl $conf 2>/dev/null`
        if [[ `echo $check` =~ "0" || -z `echo $check` ]];then
            if [[ `cat /etc/sysctl.conf` =~ "$conf" ]];then
                sed -i "s/^$conf.*/$conf=1/g" /etc/sysctl.conf
            else
                echo "$conf=1" >> /etc/sysctl.conf
            fi
            sysctl -p >/dev/null 2>&1
        fi
    done
}

main(){
    #check_sys
    offline_install
    write_service
    systemctl daemon-reload
    #set_sysctl
    systemctl enable docker.service
    systemctl start docker
    echo -e "docker install success!"
}

main
