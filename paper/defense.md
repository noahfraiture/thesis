# 1

- Hello everyone
- going to present

# 2


# 3

- 4 scenarios
- methodology
- implementation guidelines


# 4

The core of the theoritical contributions is the methodology. When I did some research online at first, it was actually pretty hard to find methodology to write cybersecurity scenarios, especially technical ones. At the beginning this thesis has a focus on the implementations of such scenarios in a cyberrange, citef, one scenarios has been implemented but the focus on the work has shift on the methodology and the scenarios themselves instead of their implementation in a particular infrastructure.

A few words about cyberranges, they are infrastructure that can support cybersecurity scenarios to train in a close and secure environment. We can think of them like a group of machine with supervisor running on top managing the machine to support VM with custom image, network and specified hardware.
And it has been shown that even about cyberrange, the literature is actually often ads disguised for a specific cyberrange. In summary the literature was of little use and this is what motivated the creation of a modular and infrastructure agnostic methodology : the helical build framework.

# 5

Here we must define a few concept that together consituate our scenarios.

# 6

First we have the scenario configuration which is the technical definition of the scenario and its infrastructure. Here we have a schema of the configuration of one of our scenario. This scenario focus on the attack of a ground base of satellite and we can see the different machine and network used in it. First the attacker will start on a machine in the internet network and will be able to access a web server exposed on internet. By gaining access to that machine with different technics not covered here for question of timing, he will enter a machine in the network of the base. We can also observe a DNS setup that will be leverage to intercept the communication with a virutal satellite.

This simple schema is used to represent the design of our system, it is not enough to define the whole scenario configuration, the system configuraiton of each machine must obviously be defined.

# 7

The attack-defense flow is maybe the most interesting part for the trainee, it is a flow representing the path the trainee will follow to achieve his goal. Here we can seee a small sample of the beginning of our scenarios about satellite. We talked that the trainee was starting on a machine connected on internet and he should access the network of the company by exploiting a web server. Here, we start with an asset, an url to the public page of the satellite company and we can start the exploit with first a virtual host discovery to discover different website on the same server behind a loadbalancer. We find another website that can be exploited with a reverse shell.
We can see that the actions are linked to a MITRE ATT&CK which was an important part of the creations of the scenarios to anchor them in a legitimity

# 8

The story line is the most simple part : it is the story of the scenario. It includes the little speech we can give the trainee to give him a goal and a motivation, but also the elements included in the scenario. Here we can see a part of the front page the trainee 

# 9

- first major contributions
- developped in parellel of scenarios
- works well with attack defense flow

# 10
# 11

# 12

Next the implementations. We will discuss two details of the implementations. It must be noted that I do not aims to define rule of best practice, only discuss some choice and give some recommendations. First the virtualization of the machines, because we cannot buy a computer for every machine of every scenarios so we must emulate them, the possibilities are containers, managed with docker or kubernetes for example, and the VM
base image for containers or vm, different but similar idea, choice of os. here only linux

# 13

Both methods have pros and cons, first let's talk about general purpose pros and cons.
A virtual machine will emulate the hardware and a container will emulate the operating system.
The virtual machine is by definition more capable than the container as it will emulate full hardware and we can run full operating system on it. It means that a kernel vulnerability like dirty cow will only be possible on virtual machine. On the other side, containers are more easy to manage and to deploy. We can use a tool like docker to build image and container and deploy them on any machine that supports docker, the management with docker compose is also very easy to add service or connect them through a network.

# 14

This was for the technology them selves, now to get something running, we need to add a base system on the virtual machine or the container. For virtual machine it will be an OS, in this master thesis I only focused on Linux distribution, so anything is possible like ubuntu fedora or even arch linux. For the container we will use a base image which can be custom to run processes, the image can start from a linux distribution, a python interpreter or pretty much any process that can run what we want. 

For both these solution we face limitations regarding the installation of our setup. For virtual machine, we need installation script in bash or ansible which are configuration file that act like evolved installation script, but during the development process in iterations, we face some issue like reentrancy : running the same script multiple time can create issues even without modifying it, revert modification is not always possible. The ultimate solution would be to delete the disk and restart an installation with the new configuration, but it is slow.
The containers are lighter and rebuild a container image is not necessarly very slow if well done so this isn't really a problem, but containers are mostly done to support a single process and can have trouble to emulate a full machine on which the attacker can enter with a shell. Also the management of dependency, so configuration file for system is not easy and the dockerfile are not made for complex file management and process configuration (talk about command that must be run when the os run ?)

# 15

For all these issue, the solution I used is NixOS. NixOS is a linux distribution with immutable design that works in a declarative manner instead of a imperative manner like traditional distribution. And all the configuration is done in a one place. This means if you want to open some ports on the firewall, you add to the configuration file of the system networking.firewall.openports = . And when you want to apply the change, you run a rebuild command that will reconfigure the whole system based on the configuration file you wrote no matter what was the previous configuration. This means a single file can define the whole configuration of a machine. From firewall configuration, to the systemd process running or the user management. And this file can obviously be copied or versioned.
A reduced version is available as a base image for container that allow us to do very similar work.

All this gives us a configuration simple to manage without complex dependency management, with a fast configuration switching.
This is the distro I personnally daily drive for almost a year.

# 16

We talked about the methodology, the implementations detail, let's see the scenarios. There is four of them
- satellite siege focus network
- hotel daemon access management of iot
- patch war secure code
- zheng hijack reverse engineer and forensic

# demo


