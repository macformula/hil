# MAC Formula HIL
Hardware-in-the-loop tester to accelerate testing of the electrical-software test-bench 

## Contributing

### Linter
- Install [golangci](https://golangci-lint.run/usage/install)
- Follow these [docs](https://golangci-lint.run/usage/integrations/) to integrate it with your IDE

## Working with Remote Linux Machines

### Requirements



- **SSH Client**: A terminal that supports SSH Client (ie. Git Bash/Powershell on Windows or Terminal on Mac). You can also use [PuTTY](https://www.putty.org/), but it is not necessary.
- **McMaster VPN**: If you plan on remoting into the machine remotely (ie. outside of Mac-WiFi), you need to download and sign into the VPN to access the network. The link for this can be found [here](https://mcmasteru365.sharepoint.com/sites/UTS-NetSoft/SitePages/All%20Software/McMasterVPN.aspx?cid=f3cfe595-8c42-48f6-a104-8baf381a39c0). Follow the installation steps and tutorial to connect to the network.

### Common Commands

#### SSH

Secure Shell (SSH) is a cryptographic network protocol used for secure communication over an unsecured network. To connect to a remote Linux machine, use the following command:

```bash
ssh username@remote_host
```
**NOTE: *remote_host* is the IP address of the machine you are accessing, and the username is unique and assigned to you.**

#### SCP

Secure Copy Protocol (SCP) allows you to securely copy files between local and remote machines. You can alternatively use **rsync** if your terminal supports it, as it transfers files faster. To copy a file from your local machine to a remote machine, use the following command:

```bash
scp local_file username@remote_host:/path/to/destination/
```
To copy a file from a remote machine to your local machine:
```bash
scp username@remote_host:/path/to/remote_file local_destination
```

### Utilities on our Raspberry Pi
#### STLink

STLink is used for programming and debugging STM32 microcontrollers.

To flash your development board with a compiled binary, run the following command. The *firmware.bin* must be built properly and sent over to the machine. 
```bash
st-flash write firmware.bin 0x08000000
```

#### SocketCAN
SocketCAN provides a set of kernel modules and utilities for CAN networking on Linux. CanUtils includes command-line tools for interacting with CAN networks. Use the tools for various tasks, such as sending and receiving CAN messages:
```bash
# Example command to send a CAN message on CAN0
cansend can0 123#1122334455667788

# Example command to receive CAN messages on CAN0
candump can0
```