package iosxe

const hwFilter = `<device-hardware-data xmlns="http://cisco.com/ns/yang/Cisco-IOS-XE-device-hardware-oper">
    <device-hardware>
      <device-inventory/>
    </device-hardware>
</device-hardware-data>`

const systemFilter = `<system xmlns="http://openconfig.net/yang/system">
 <config>
 </config>
  <state>
	  <hostname/>
		<domain-name/>
  </state>
</system>
`

const interfaceFilter = `<interfaces xmlns="http://openconfig.net/yang/interfaces">
    <interface>
      <name/>
      <state>
        <description/>
        <name/>
        <type/>
        <enabled/>
      </state>
      <ethernet xmlns="http://openconfig.net/yang/interfaces/ethernet">
        <state>
          <mac-address/>
          <auto-negotiate/>
          <port-speed/>
        </state>
      </ethernet>
    </interface>
  </interfaces>`

const arpFilter = `<arp-data xmlns="http://cisco.com/ns/yang/Cisco-IOS-XE-arp-oper"/>`
