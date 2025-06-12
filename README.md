# gvswitch

This is a simple utility for running completely userspace low-privileged virtual
L2+ switch with additional features like simple routing and DHCP. It is useful
for QEMU virtual machines without root permissions (especially on MacOS).

The implementation is mostly based on gvisor implementation of network stack.

## Usage

Foremost, you need a configuration file. Although it has no default config file,
it has some routines arount that. You can generate it, show (with all computed
values), and validate. Some of the fields might be filled out automatically,
e.g. gateway IP and hypervisor IP will be the first and the last valid IP
addresses correspondingly.

```bash
###
### Generate the config template with per-field explanation
###
gvswitch config template

###
### Show config, with default/computed values
###
gvswitch config show -c config.yaml

###
### Validate config file
###
gvswitch validate -c config.yaml
```

## Connecting VMs to the switch

### QEMU

As any virtual machine, it's also very common to enable networking for QEMU too.
You can connect your QEMU virtual machine to the virtual switch with the following
arguments. Please keep in mind that normally there are hundreds of such arguments
for real VMs, and this short instruction only covers the arguments for connection
to the virtual switch.

```bash
qemu-system-aarch64 \
    -netdev '{"type":"socket","connect":"127.0.0.1:58557","id":"hostnet0"}' \
    -device '{"driver":"virtio-net-pci","netdev":"hostnet0","id":"net0","mac":"80:e2:12:ee:81:bb"}'
```

### Libvirt

You can also use it with libvirt hypervisor. Actually, libvirt is a kind of
wrapper around QEMU virtual machine, thus it's sequential that it is also
supported.

In the domain (VM) definition, use the following or similar configuration:

```xml
<domain ...>
  ...
  <devices>
    ...
    <!-- gvswitch hosts the server side -->
    <interface type='client'>
      <!--
        optionally, you can define MAC address and set a static
        IP address for the particular virtual machine
      -->
      <mac address='80:e2:12:ee:81:bb'/>
      <!--
        This is the port listened by gvswitch. Not the IP address
        of VM. In the gvswitch config file, it should be in "serve.qemu".
      -->
      <source address='127.0.0.1' port='58557'/>
    </interface>
...
```

## Notes about the project

This project is a small part of another educational project, known as
"Kubernetes - The Mindful Way". This virtual switch used for prototyping
and testing the related Ansible and Molecule artifacts. Of course, it is
also mentioned in the learning materials.

Originally, it is an [idea of RedHat engineers](https://github.com/containers/gvisor-tap-vsock). However, their solution
doesn't allow applying it with "Kubernetes - The Mindful Way" project,
thus it is just a reimplementation with more flexible and suitable
configuration interface.

"Kubernetes - The Mindful Way" is a public project which will be available
pretty soon (if not already).
