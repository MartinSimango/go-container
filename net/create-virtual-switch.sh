#./delete-veth.sh || true

# echo "deleting inteface"

# sudo ip link del v-net-0
# sudo ip -n red link del veth-red
# sudo ip -n blue link del veth-blue

# create virtual switch

echo "creating virtual switch"

sudo ip link add mydocker type bridge
sudo ip link set mydocker up
sudo ip addr add 172.17.0.1/16 dev mydocker
# create veth pairs

echo "creating veth pairs"
sudo ip link add veth-cont type veth peer name veth-cont-br

sudo ip link add veth-blue type veth peer name veth-blue-br


# attach interfaces to namespaces and to switch

echo "attaching virtual interfaces to namespaces and switch"
sudo ip link set veth-cont netns 4666 # pid of container (which has it's on network namespace) 
sudo ip link set veth-cont-br master mydocker

sudo ip link set veth-blue netns blue 
sudo ip link set veth-blue-br master v-net-0  

# attach ip addresses

echo "assinging ip address and setting interfaces up"
sudo ip -n red addr add 192.168.15.1/24 dev veth-red

sudo ip addr 172.17.0.2/16 dev veth-cont


sudo ip -n blue addr add 192.168.15.2/24 dev veth-blue

sudo ip -n red link set veth-red up

sudo ip -n blue link set veth-blue up
sudo ip link set veth-red-br up
sudo ip link set veth-blue-br up


# adding routing tables

echo "setting up routing tables"

sudo ip -n red route add default via 192.168.15.3
sudo ip -n blue route add default via 192.168.15.3


iptables -t nat -A POSTROUTING -s 172.17.0.0/16 -j MASQUERADE