#delete
sudo ip -n red link del veth-red

#create network namespaces
#sudo ip netns add red
#sudo ip netns add blue

#create veth pear
sudo ip link add veth-red type veth peer name veth-blue

#attach virtual interfaces to namespaces
sudo ip link set veth-red netns red
sudo ip link set veth-blue netns blue

#assign ip addresses to namespaces
sudo ip -n red addr add 192.168.15.1/24 dev veth-red
sudo ip -n blue addr add 192.168.15.2/24 dev veth-blue

#bring up the interfaces

sudo ip -n red link set veth-red up
sudo ip -n blue link set veth-blue up
sudo ip link set veth-red-br up
sudo ip link set veth-blue-br up