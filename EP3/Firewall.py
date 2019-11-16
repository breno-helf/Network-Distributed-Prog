# Copyright 2012 James McCauley
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at:
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""
This component is for use with the OpenFlow tutorial.

It acts as a simple hub, but can be modified to act like an L2
learning switch.

It's roughly similar to the one Brandon Heller did for NOX.
"""

from pox.core import core
import pox.lib.packet as pkt
import pox.openflow.libopenflow_01 as of
from pox.lib.addresses import IPAddr, IPAddr6, EthAddr

log = core.getLogger()

SRC_IP_ADDRESS = None #IPAddr("10.0.0.2")
DST_IP_ADDRESS = None #IPAddr("10.0.0.2")
PORT = None #5001
PROTOCOL_TYPE = None #"udp"  

class Tutorial (object):
  """
  A Tutorial object is created for each switch that connects.
  A Connection object for that switch is passed to the __init__ function.
  """
  def __init__ (self, connection):
    # Keep track of the connection to the switch so that we can
    # send it messages!
    self.connection = connection

    # This binds our PacketIn event listener
    connection.addListeners(self)

    # Use this table to keep track of which ethernet address is on
    # which switch port (keys are MACs, values are ports).
    self.mac_to_port = {}
    self.iparp = {}


  def resend_packet (self, packet_in, out_port):
    """
    Instructs the switch to resend a packet that it had sent to us.
    "packet_in" is the ofp_packet_in object the switch had sent to the
    controller due to a table-miss.
    """
    msg = of.ofp_packet_out()
    msg.data = packet_in

    # Add an action to send to the specified port
    action = of.ofp_action_output(port = out_port)
    msg.actions.append(action)

    # Send message to switch
    self.connection.send(msg)


  def act_like_hub (self, packet, packet_in):
    """
    Implement hub-like behavior -- send all packets to all ports besides
    the input port.
    """

    # We want to output to all ports -- we do that using the special
    # OFPP_ALL port as the output port.  (We could have also used
    # OFPP_FLOOD.)
    self.resend_packet(packet_in, of.OFPP_ALL)

    # Note that if we didn't get a valid buffer_id, a slightly better
    # implementation would check that we got the full data before
    # sending it (len(packet_in.data) should be == packet_in.total_len)).


  def act_like_switch (self, packet, packet_in):

    self.mac_to_port[packet.src] = packet_in.in_port

    # if(str(packet.dst) == "ff:ff:ff:ff:ff:ff"):
    #   if packet.src in self.mac_to_port:
    #     if(self.is_mac_ip(SRC_IP_ADDRESS, packet.src)):
    #       return
    #     else:
    #       self.resend_packet()
    # print("DST: " + str(packet.dst))
    # print(packet)

    if packet.dst in self.mac_to_port:

      # log.debug("Installing flow. Source: " + str(packet.src) + ". Destination: " + str(packet.dst) + ". Port: "+ str(self.mac_to_port[packet.dst]))

      msg = of.ofp_flow_mod()
      msg.match = of.ofp_match.from_packet(packet)
      msg.idle_timeout = 10
      msg.buffer_id = packet_in.buffer_id
      msg.actions.append(of.ofp_action_output(port = self.mac_to_port[packet.dst]))
      self.connection.send(msg)

    else:
      # print (str(packet.dst) + " not known, resend to everybody")
      self.resend_packet(packet_in, of.OFPP_ALL)

  def handle_IP_packet (self, packet):
    ip = packet.find('ipv4')
    if ip is None:
      # This packet isn't IP!
      return 0

    # self.iparp[ip.srcip] = packet.src
    # self.iparp[ip.dstip] = packet.dst
    if SRC_IP_ADDRESS != None:
      if ip.srcip == SRC_IP_ADDRESS:
        return 1
    if DST_IP_ADDRESS != None:
      if ip.dstip == DST_IP_ADDRESS:
        return 1
    print ("Source IP:"), ip.srcip
    print ("Destin IP:"), ip.dstip
    return 0

  def check_port(self, packet):
    if PORT == None:
      return 0
    match = of.ofp_match.from_packet(packet)
    if match.tp_dst == PORT :
      print("PORT_FILTER")
      return 1
    if match.tp_src == PORT :
      print("PORT_FILTER")
      return 1
    return 0  

  def Protocol_Type (self, packet):
    if PROTOCOL_TYPE == None:
      return 0
    ip = packet.find(PROTOCOL_TYPE)
    if ip is None:
      # This packet isn't IP!
      return 0 
    return 1

  def _handle_PacketIn (self, event):
    """
    Handles packet in messages from the switch.
    """

    packet = event.parsed # This is the parsed packet data.


    if not packet.parsed:
      log.warning("Ignoring incomplete packet")
      return

    packet_in = event.ofp # The actual ofp_packet_in message.

    # msg = of.ofp_packet_out()
    # msg.data = packet_in
    # print(msg)
    # print("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
    # Comment out the following line and uncomment the one after
    # when starting the exercise.
    # self.act_like_hub(packet, packet_in)
    self.check_port(packet)
    if self.Protocol_Type(packet) == 1:
      print("PROTOCOL_FILTER")
    if self.handle_IP_packet(packet) == 0:
      # if self.is_mac_ip(SRC_IP_ADDRESS, packet.src) == 0:
      #   if self.is_mac_ip(DST_IP_ADDRESS, packet.dst) == 0:
      #     print("PASSOU")
      #     print ("SRC IP:"), packet.src
      #     print ("DST IP:"), packet.dst 
      self.act_like_switch(packet, packet_in)

    
    match = of.ofp_match.from_packet(packet)
    print (match.show())
    print("------------------------------------------------------")



def launch ():
  """
  Starts the component
  """
  def start_switch (event):
    log.debug("Controlling %s" % (event.connection,))
    Tutorial(event.connection)
  core.openflow.addListenerByName("ConnectionUp", start_switch)
