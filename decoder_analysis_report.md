# GoPacket Decoder Analysis Report

## Overview

Total decoders analyzed: 23
Total issues found: 851

## arp.go

### Bounds Checking

- Potential bounds issue: data[2:4] without prior length check
- Potential bounds issue: data[4] without prior length check
- Potential bounds issue: data[5] without prior length check
- Potential bounds issue: data[6:8] without prior length check
- Potential bounds issue: bytes[2:] without prior length check
- Potential bounds issue: bytes[4] without prior length check
- Potential bounds issue: bytes[5] without prior length check
- Potential bounds issue: bytes[6:] without prior length check

### Rfc Compliance

- Magic number 16 should be a named constant
- Magic number 2009 should be a named constant
- Magic number 2012 should be a named constant

### Security

- Potential integer overflow: arp * ARP
- Potential integer overflow: arp * ARP
- Potential integer overflow: arp * ARP
- Potential integer overflow: arp * ARP
- Potential integer overflow: arp * ARP

## ethernet.go

### Bounds Checking

- Potential bounds issue: data[12:14] without prior length check
- Potential bounds issue: data[14:] without prior length check
- Potential bounds issue: bytes[6:] without prior length check
- Potential bounds issue: bytes[12:] without prior length check
- Potential bounds issue: bytes[12:] without prior length check

### Rfc Compliance

- Magic number 14 should be a named constant
- Magic number 12 should be a named constant
- Magic number 2009 should be a named constant
- Magic number 0600 should be a named constant
- Magic number 2012 should be a named constant
- Magic number 16 should be a named constant
- Magic number 60 should be a named constant
- Magic number 802 should be a named constant

### Security

- Potential integer overflow: e * Ethernet
- Potential integer overflow: e * Ethernet
- Potential integer overflow: eth * Ethernet
- Potential integer overflow: eth * Ethernet
- Potential integer overflow: eth * Ethernet
- Potential integer overflow: eth * Ethernet

## ip4.go

### Bounds Checking

- Potential bounds issue: bytes[0] without prior length check
- Potential bounds issue: bytes[1] without prior length check
- Potential bounds issue: bytes[2:] without prior length check
- Potential bounds issue: bytes[4:] without prior length check
- Potential bounds issue: bytes[6:] without prior length check
- Potential bounds issue: bytes[8] without prior length check
- Potential bounds issue: bytes[9] without prior length check
- Potential bounds issue: bytes[12:16] without prior length check
- Potential bounds issue: bytes[16:20] without prior length check
- Potential bounds issue: bytes[10:] without prior length check
- Potential bounds issue: bytes[10] without prior length check
- Potential bounds issue: bytes[11] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1] without prior length check
- Potential bounds issue: data[2:4] without prior length check
- Potential bounds issue: data[4:6] without prior length check
- Potential bounds issue: data[8] without prior length check
- Potential bounds issue: data[9] without prior length check
- Potential bounds issue: data[10:12] without prior length check
- Potential bounds issue: data[12:16] without prior length check
- Potential bounds issue: data[16:20] without prior length check
- Potential bounds issue: data[1:] without prior length check
- Potential bounds issue: data[1:] without prior length check

### Field Validation

- Missing version field validation
- Flag field I defined but not validated

### Rfc Compliance

- Magic number 32 should be a named constant
- Magic number 13 should be a named constant
- Magic number 12 should be a named constant
- Magic number 2009 should be a named constant
- Magic number 3514 should be a named constant
- Magic number 16 should be a named constant
- Magic number 2012 should be a named constant
- Magic number 11 should be a named constant

### Security

- Potential integer overflow: i * IPv4
- Potential integer overflow: i * IPv4
- Potential integer overflow: ip * IPv4
- Potential integer overflow: ip * IPv4
- Potential integer overflow: ip * IPv4
- Potential integer overflow: ip * IPv4
- Potential integer overflow: i * IPv4
- Potential integer overflow: i * IPv4
- Potential integer overflow: ip * IPv4
- Potential infinite loop detected
- Potential infinite loop detected
- Unchecked allocation size: 4
- Unchecked allocation size: 0

## ip6.go

### Bounds Checking

- Potential bounds issue: hbh[1] without prior length check
- Potential bounds issue: bytes[0] without prior length check
- Potential bounds issue: bytes[1] without prior length check
- Potential bounds issue: bytes[2:] without prior length check
- Potential bounds issue: bytes[4:] without prior length check
- Potential bounds issue: bytes[6] without prior length check
- Potential bounds issue: bytes[7] without prior length check
- Potential bounds issue: bytes[8:] without prior length check
- Potential bounds issue: bytes[24:] without prior length check
- Potential bounds issue: data[0:2] without prior length check
- Potential bounds issue: data[0:4] without prior length check
- Potential bounds issue: data[4:6] without prior length check
- Potential bounds issue: data[6] without prior length check
- Potential bounds issue: data[7] without prior length check
- Potential bounds issue: data[8:24] without prior length check
- Potential bounds issue: data[24:40] without prior length check
- Potential bounds issue: data[40:] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1] without prior length check
- Potential bounds issue: data[2:] without prior length check
- Potential bounds issue: data[1] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1] without prior length check
- Potential bounds issue: data[2:] without prior length check
- Potential bounds issue: OptionAlignment[0] without prior length check
- Potential bounds issue: OptionAlignment[1] without prior length check
- Potential bounds issue: bytes[0] without prior length check
- Potential bounds issue: bytes[1] without prior length check
- Potential bounds issue: data[2] without prior length check
- Potential bounds issue: data[3] without prior length check
- Potential bounds issue: data[4:8] without prior length check
- Potential bounds issue: Contents[8:] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1] without prior length check
- Potential bounds issue: data[2:4] without prior length check
- Potential bounds issue: data[3] without prior length check
- Potential bounds issue: data[3] without prior length check
- Potential bounds issue: data[4:8] without prior length check
- Potential bounds issue: bytes[0] without prior length check
- Potential bounds issue: bytes[1] without prior length check
- DecodeFromBytes function doesn't check data length
- Missing minimum length validation in DecodeFromBytes

### Field Validation

- Missing version field validation
- Reserved fields not validated to be zero

### Rfc Compliance

- Magic number 32 should be a named constant
- Magic number 2009 should be a named constant
- Magic number 2012 should be a named constant
- Magic number 29 should be a named constant
- Magic number 24 should be a named constant
- Magic number 2460 should be a named constant
- Magic number 2675 should be a named constant
- Magic number 40 should be a named constant

### Security

- Potential integer overflow: HopByHop * IPv6HopByHop
- Potential integer overflow: ipv6 * IPv6
- Potential integer overflow: ipv6 * IPv6
- Potential integer overflow: hopopts * IPv6HopByHop
- Potential integer overflow: tlv * IPv6HopByHopOption
- Potential integer overflow: ip6 * IPv6
- Potential integer overflow: tlv * IPv6HopByHopOption
- Potential integer overflow: ipv6 * IPv6
- Potential integer overflow: ipv6 * IPv6
- Potential integer overflow: ipv6 * IPv6
- Potential integer overflow: ipv6 * IPv6
- Potential integer overflow: h * ipv6HeaderTLVOption
- Potential integer overflow: h * ipv6HeaderTLVOption
- Potential integer overflow: x * n
- Potential integer overflow: i * IPv6ExtensionSkipper
- Potential integer overflow: i * IPv6ExtensionSkipper
- Potential integer overflow: i * IPv6ExtensionSkipper
- Potential integer overflow: i * IPv6HopByHop
- Potential integer overflow: i * IPv6HopByHop
- Potential integer overflow: i * IPv6HopByHop
- Potential integer overflow: o * IPv6HopByHopOption
- Potential integer overflow: i * IPv6Routing
- Potential integer overflow: i * IPv6Fragment
- Potential integer overflow: i * IPv6Destination
- Potential integer overflow: i * IPv6Destination
- Potential integer overflow: i * IPv6Destination
- Potential integer overflow: ipv6 * IPv6
- Potential infinite loop detected
- Unchecked allocation size: 4

## tcp.go

### Bounds Checking

- Potential bounds issue: OptionData[4:8] without prior length check
- Potential bounds issue: bytes[2:] without prior length check
- Potential bounds issue: bytes[4:] without prior length check
- Potential bounds issue: bytes[8:] without prior length check
- Potential bounds issue: bytes[12:] without prior length check
- Potential bounds issue: bytes[14:] without prior length check
- Potential bounds issue: bytes[18:] without prior length check
- Potential bounds issue: bytes[16] without prior length check
- Potential bounds issue: bytes[17] without prior length check
- Potential bounds issue: bytes[16:] without prior length check
- Potential bounds issue: data[2:4] without prior length check
- Potential bounds issue: data[2:4] without prior length check
- Potential bounds issue: data[4:8] without prior length check
- Potential bounds issue: data[8:12] without prior length check
- Potential bounds issue: data[12] without prior length check
- Potential bounds issue: data[13] without prior length check
- Potential bounds issue: data[13] without prior length check
- Potential bounds issue: data[13] without prior length check
- Potential bounds issue: data[13] without prior length check
- Potential bounds issue: data[13] without prior length check
- Potential bounds issue: data[13] without prior length check
- Potential bounds issue: data[13] without prior length check
- Potential bounds issue: data[13] without prior length check
- Potential bounds issue: data[12] without prior length check
- Potential bounds issue: data[14:16] without prior length check
- Potential bounds issue: data[16:18] without prior length check
- Potential bounds issue: data[18:20] without prior length check
- Potential bounds issue: data[1:] without prior length check

### Rfc Compliance

- Magic number 0040 should be a named constant
- Magic number 80 should be a named constant
- Magic number 2012 should be a named constant
- Magic number 16 should be a named constant
- Magic number 0020 should be a named constant
- Magic number 40 should be a named constant
- Magic number 32 should be a named constant
- Magic number 0100 should be a named constant
- Magic number 17 should be a named constant
- Magic number 18 should be a named constant
- Magic number 2009 should be a named constant
- Magic number 0080 should be a named constant

### Security

- Potential integer overflow: t * TCP
- Potential integer overflow: t * TCP
- Potential integer overflow: t * TCP
- Potential integer overflow: t * TCP
- Potential integer overflow: tcp * TCP
- Potential integer overflow: t * TCP
- Potential integer overflow: t * TCP
- Potential integer overflow: t * TCP
- Potential integer overflow: t * TCP
- Unchecked allocation size: 2
- Unchecked allocation size: 2

## udp.go

### Bounds Checking

- Potential bounds issue: data[2:4] without prior length check
- Potential bounds issue: data[2:4] without prior length check
- Potential bounds issue: data[4:6] without prior length check
- Potential bounds issue: data[6:8] without prior length check
- Potential bounds issue: bytes[2:] without prior length check
- Potential bounds issue: bytes[4:] without prior length check
- Potential bounds issue: bytes[6] without prior length check
- Potential bounds issue: bytes[7] without prior length check
- Potential bounds issue: bytes[6:] without prior length check

### Rfc Compliance

- Magic number 16 should be a named constant
- Magic number 65535 should be a named constant
- Magic number 2009 should be a named constant
- Magic number 2012 should be a named constant

### Security

- Potential integer overflow: u * UDP
- Potential integer overflow: udp * UDP
- Potential integer overflow: u * UDP
- Potential integer overflow: u * UDP
- Potential integer overflow: u * UDP
- Potential integer overflow: u * UDP
- Potential integer overflow: u * UDP
- Unchecked allocation size: 2
- Unchecked allocation size: 2

## icmp4.go

### Bounds Checking

- Potential bounds issue: data[4:6] without prior length check
- Potential bounds issue: data[6:8] without prior length check
- Potential bounds issue: data[8:] without prior length check
- Potential bounds issue: bytes[4:] without prior length check
- Potential bounds issue: bytes[6:] without prior length check
- Potential bounds issue: bytes[2] without prior length check
- Potential bounds issue: bytes[3] without prior length check
- Potential bounds issue: bytes[2:] without prior length check

### Rfc Compliance

- Magic number 2009 should be a named constant
- Magic number 2012 should be a named constant

### Security

- Potential integer overflow: codeStr * map
- Potential integer overflow: i * ICMPv4
- Potential integer overflow: i * ICMPv4
- Potential integer overflow: i * ICMPv4
- Potential integer overflow: i * ICMPv4
- Potential integer overflow: i * ICMPv4

## icmp6.go

### Bounds Checking

- Potential bounds issue: data[4:] without prior length check
- Potential bounds issue: bytes[2] without prior length check
- Potential bounds issue: bytes[3] without prior length check
- Potential bounds issue: bytes[2:] without prior length check

### Rfc Compliance

- Magic number 4443 should be a named constant
- Magic number 3810 should be a named constant
- Magic number 2710 should be a named constant
- Magic number 2012 should be a named constant
- Magic number 16 should be a named constant
- Magic number 20 should be a named constant
- Magic number 2009 should be a named constant
- Magic number 4861 should be a named constant

### Security

- Potential integer overflow: codeStr * map
- Potential integer overflow: i * ICMPv6
- Potential integer overflow: i * ICMPv6
- Potential integer overflow: i * ICMPv6
- Potential integer overflow: i * ICMPv6
- Potential integer overflow: i * ICMPv6

## dhcpv4.go

### Bounds Checking

- Potential bounds issue: data[2] without prior length check
- Potential bounds issue: data[3] without prior length check
- Potential bounds issue: data[4:8] without prior length check
- Potential bounds issue: data[8:10] without prior length check
- Potential bounds issue: data[10:12] without prior length check
- Potential bounds issue: data[12:16] without prior length check
- Potential bounds issue: data[16:20] without prior length check
- Potential bounds issue: data[20:24] without prior length check
- Potential bounds issue: data[24:28] without prior length check
- Potential bounds issue: data[44:108] without prior length check
- Potential bounds issue: data[108:236] without prior length check
- Potential bounds issue: data[236:240] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1] without prior length check
- Potential bounds issue: data[2] without prior length check
- Potential bounds issue: data[3] without prior length check
- Potential bounds issue: data[4:8] without prior length check
- Potential bounds issue: data[8:10] without prior length check
- Potential bounds issue: data[10:12] without prior length check
- Potential bounds issue: data[12:16] without prior length check
- Potential bounds issue: data[16:20] without prior length check
- Potential bounds issue: data[20:24] without prior length check
- Potential bounds issue: data[24:28] without prior length check
- Potential bounds issue: data[28:44] without prior length check
- Potential bounds issue: data[44:108] without prior length check
- Potential bounds issue: data[108:236] without prior length check
- Potential bounds issue: data[236:240] without prior length check
- Potential bounds issue: Data[0] without prior length check
- Potential bounds issue: Data[0] without prior length check
- Potential bounds issue: Data[1] without prior length check
- Potential bounds issue: Data[2] without prior length check
- Potential bounds issue: Data[3] without prior length check
- Potential bounds issue: b[0] without prior length check
- Potential bounds issue: b[0] without prior length check
- Potential bounds issue: b[1] without prior length check
- Potential bounds issue: b[2:] without prior length check

### Rfc Compliance

- Magic number 63825363 should be a named constant
- Magic number 2132 should be a named constant
- Magic number 2016 should be a named constant
- Magic number 2131 should be a named constant
- Magic number 868 should be a named constant
- Magic number 108 should be a named constant
- Magic number 236 should be a named constant
- Magic number 116 should be a named constant

### Security

- Potential integer overflow: d * DHCPv4
- Potential integer overflow: d * DHCPv4
- Potential integer overflow: d * DHCPv4
- Potential integer overflow: d * DHCPv4
- Potential integer overflow: d * DHCPv4
- Potential integer overflow: d * DHCPv4
- Potential integer overflow: o * DHCPOption
- Potential integer overflow: o * DHCPOption

## dhcpv6.go

### Bounds Checking

- Potential bounds issue: data[18:34] without prior length check
- Potential bounds issue: data[1:4] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1] without prior length check
- Potential bounds issue: data[2:18] without prior length check
- Potential bounds issue: data[18:34] without prior length check
- Potential bounds issue: data[1:4] without prior length check
- Potential bounds issue: data[0:2] without prior length check
- Potential bounds issue: data[2:4] without prior length check
- Potential bounds issue: data[4:8] without prior length check
- Potential bounds issue: data[8:] without prior length check
- Potential bounds issue: data[2:6] without prior length check
- Potential bounds issue: data[6:] without prior length check
- Potential bounds issue: data[4:] without prior length check

### Rfc Compliance

- Magic number 18 should be a named constant
- Magic number 19 should be a named constant
- Magic number 16 should be a named constant
- Magic number 3315 should be a named constant
- Magic number 2018 should be a named constant

### Security

- Potential integer overflow: d * DHCPv6
- Potential integer overflow: d * DHCPv6
- Potential integer overflow: d * DHCPv6
- Potential integer overflow: d * DHCPv6
- Potential integer overflow: d * DHCPv6
- Potential integer overflow: d * DHCPv6
- Potential integer overflow: d * DHCPv6DUID
- Potential integer overflow: d * DHCPv6DUID
- Potential integer overflow: d * DHCPv6DUID
- Potential integer overflow: d * DHCPv6DUID
- Unchecked allocation size: length

## dns.go

### Bounds Checking

- Potential bounds issue: data[2] without prior length check
- Potential bounds issue: data[2] without prior length check
- Potential bounds issue: data[2] without prior length check
- Potential bounds issue: data[2] without prior length check
- Potential bounds issue: data[2] without prior length check
- Potential bounds issue: data[3] without prior length check
- Potential bounds issue: data[3] without prior length check
- Potential bounds issue: data[3] without prior length check
- Potential bounds issue: data[4:6] without prior length check
- Potential bounds issue: data[6:8] without prior length check
- Potential bounds issue: data[8:10] without prior length check
- Potential bounds issue: data[10:12] without prior length check
- Potential bounds issue: bytes[2] without prior length check
- Potential bounds issue: bytes[3] without prior length check
- Potential bounds issue: bytes[4:] without prior length check
- Potential bounds issue: bytes[6:] without prior length check
- Potential bounds issue: bytes[8:] without prior length check
- Potential bounds issue: bytes[10:] without prior length check
- Potential bounds issue: Data[4:] without prior length check

### Field Validation

- Missing version field validation
- Reserved fields not validated to be zero

### Rfc Compliance

- Magic number 2782 should be a named constant
- Magic number 7553 should be a named constant
- Magic number 63 should be a named constant
- Magic number 2671 should be a named constant
- Magic number 2018 should be a named constant
- Magic number 80 should be a named constant
- Magic number 3596 should be a named constant
- Magic number 4635 should be a named constant
- Magic number 6195 should be a named constant
- Magic number 40 should be a named constant
- Magic number 32 should be a named constant
- Magic number 87 should be a named constant
- Magic number 2014 should be a named constant
- Magic number 1035 should be a named constant
- Magic number 2673 should be a named constant
- Magic number 2845 should be a named constant
- Magic number 3425 should be a named constant
- Magic number 2136 should be a named constant
- Magic number 7873 should be a named constant
- Magic number 1034 should be a named constant
- Magic number 1996 should be a named constant
- Magic number 6891 should be a named constant
- Magic number 2930 should be a named constant

### Security

- Potential integer overflow: d * DNS
- Potential integer overflow: d * DNS
- Potential integer overflow: d * DNS
- Potential integer overflow: d * DNS
- Potential integer overflow: d * DNS
- Potential integer overflow: rr * DNSResourceRecord
- Potential integer overflow: d * DNS
- Potential integer overflow: q * DNSQuestion
- Potential integer overflow: q * DNSQuestion
- Potential integer overflow: rr * DNSResourceRecord
- Potential integer overflow: rr * DNSResourceRecord
- Potential integer overflow: rr * DNSResourceRecord
- Potential integer overflow: rr * DNSResourceRecord
- Potential infinite loop detected
- Unchecked allocation size: 1
- Unchecked allocation size: 0

## igmp.go

### Bounds Checking

- Potential bounds issue: data[8] without prior length check
- Potential bounds issue: data[4:8] without prior length check
- Potential bounds issue: data[8] without prior length check
- Potential bounds issue: data[9] without prior length check
- Potential bounds issue: data[10:12] without prior length check
- Potential bounds issue: data[4:8] without prior length check
- Potential bounds issue: byte[0] without prior length check

### Field Validation

- Missing version field validation
- Reserved fields not validated to be zero

### Rfc Compliance

- Magic number 80 should be a named constant
- Magic number 2012 should be a named constant
- Magic number 16 should be a named constant
- Magic number 70 should be a named constant
- Magic number 32 should be a named constant
- Magic number 17 should be a named constant
- Magic number 100 should be a named constant
- Magic number 22 should be a named constant
- Magic number 2009 should be a named constant
- Magic number 3376 should be a named constant
- Magic number 11 should be a named constant

### Security

- Potential integer overflow: i * IGMPv1or2
- Potential integer overflow: i * IGMP
- Potential integer overflow: i * IGMP
- Potential integer overflow: i * IGMP
- Potential integer overflow: i * IGMPv1or2
- Potential integer overflow: i * IGMPv1or2
- Potential integer overflow: i * IGMPv1or2
- Potential integer overflow: i * IGMPv1or2
- Potential integer overflow: i * IGMP
- Potential integer overflow: i * IGMP
- Potential integer overflow: i * IGMP
- Potential infinite loop detected

## ospf.go

### Bounds Checking

- Potential bounds issue: data[20:24] without prior length check
- Potential bounds issue: data[24] without prior length check
- Potential bounds issue: data[24:28] without prior length check
- Potential bounds issue: data[28:32] without prior length check
- Potential bounds issue: data[32:36] without prior length check
- Potential bounds issue: data[20:24] without prior length check
- Potential bounds issue: data[20] without prior length check
- Potential bounds issue: data[20:24] without prior length check
- Potential bounds issue: data[20:24] without prior length check
- Potential bounds issue: data[20:24] without prior length check
- Potential bounds issue: data[24] without prior length check
- Potential bounds issue: data[25] without prior length check
- Potential bounds issue: data[20:24] without prior length check
- Potential bounds issue: data[24:28] without prior length check
- Potential bounds issue: data[28:32] without prior length check
- Potential bounds issue: data[20] without prior length check
- Potential bounds issue: data[24] without prior length check
- Potential bounds issue: data[20:24] without prior length check
- Potential bounds issue: data[25] without prior length check
- Potential bounds issue: data[26:28] without prior length check
- Potential bounds issue: data[40:44] without prior length check
- Potential bounds issue: data[20] without prior length check
- Potential bounds issue: data[20:24] without prior length check
- Potential bounds issue: data[24:40] without prior length check
- Potential bounds issue: data[20:22] without prior length check
- Potential bounds issue: data[22:24] without prior length check
- Potential bounds issue: data[24:28] without prior length check
- Potential bounds issue: data[28:32] without prior length check
- Potential bounds issue: data[4:8] without prior length check
- Potential bounds issue: data[8:12] without prior length check
- Potential bounds issue: data[12:14] without prior length check
- Potential bounds issue: data[14:16] without prior length check
- Potential bounds issue: data[16:24] without prior length check
- Potential bounds issue: data[24:28] without prior length check
- Potential bounds issue: data[28:30] without prior length check
- Potential bounds issue: data[30] without prior length check
- Potential bounds issue: data[31] without prior length check
- Potential bounds issue: data[32:36] without prior length check
- Potential bounds issue: data[36:40] without prior length check
- Potential bounds issue: data[40:44] without prior length check
- Potential bounds issue: data[24:26] without prior length check
- Potential bounds issue: data[26] without prior length check
- Potential bounds issue: data[27] without prior length check
- Potential bounds issue: data[28:32] without prior length check
- Potential bounds issue: data[24:28] without prior length check
- Potential bounds issue: data[28:] without prior length check
- Potential bounds issue: data[4:8] without prior length check
- Potential bounds issue: data[8:12] without prior length check
- Potential bounds issue: data[12:14] without prior length check
- Potential bounds issue: data[14] without prior length check
- Potential bounds issue: data[15] without prior length check
- Potential bounds issue: data[16:20] without prior length check
- Potential bounds issue: data[20] without prior length check
- Potential bounds issue: data[21:25] without prior length check
- Potential bounds issue: data[24:26] without prior length check
- Potential bounds issue: data[26:28] without prior length check
- Potential bounds issue: data[28:32] without prior length check
- Potential bounds issue: data[32:36] without prior length check
- Potential bounds issue: data[16:20] without prior length check
- Potential bounds issue: data[20:22] without prior length check
- Potential bounds issue: data[22:24] without prior length check
- Potential bounds issue: data[24:28] without prior length check
- Potential bounds issue: data[16:20] without prior length check
- Potential bounds issue: data[20:] without prior length check

### Field Validation

- Missing version field validation
- Reserved fields not validated to be zero

### Rfc Compliance

- Magic number 14 should be a named constant
- Magic number 30 should be a named constant
- Magic number 31 should be a named constant
- Magic number 2002 should be a named constant
- Magic number 2017 should be a named constant
- Magic number 2004 should be a named constant
- Magic number 2003 should be a named constant
- Magic number 80 should be a named constant
- Magic number 2328 should be a named constant
- Magic number 4005 should be a named constant
- Magic number 2001 should be a named constant
- Magic number 64 should be a named constant
- Magic number 2007 should be a named constant
- Magic number 15 should be a named constant
- Magic number 26 should be a named constant
- Magic number 40 should be a named constant
- Magic number 5340 should be a named constant
- Magic number 25 should be a named constant
- Magic number 22 should be a named constant
- Magic number 18 should be a named constant
- Magic number 21 should be a named constant
- Magic number 2009 should be a named constant
- Magic number 27 should be a named constant

### Security

- Potential integer overflow: ospf * OSPFv2
- Potential integer overflow: ospf * OSPFv3
- Potential integer overflow: ospf * OSPFv2
- Potential integer overflow: ospf * OSPFv3
- Potential integer overflow: ospf * OSPFv2
- Potential integer overflow: ospf * OSPFv3
- Potential integer overflow: ospf * OSPFv2
- Potential integer overflow: ospf * OSPFv3
- Potential infinite loop detected

## ntp.go

### Bounds Checking

- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1] without prior length check
- Potential bounds issue: data[2] without prior length check
- Potential bounds issue: data[3] without prior length check
- Potential bounds issue: data[4:8] without prior length check
- Potential bounds issue: data[8:12] without prior length check
- Potential bounds issue: data[12:16] without prior length check
- Potential bounds issue: data[16:24] without prior length check
- Potential bounds issue: data[24:32] without prior length check
- Potential bounds issue: data[32:40] without prior length check
- Potential bounds issue: data[40:48] without prior length check
- Potential bounds issue: data[48:] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1] without prior length check
- Potential bounds issue: data[2] without prior length check
- Potential bounds issue: data[3] without prior length check
- Potential bounds issue: data[4:8] without prior length check
- Potential bounds issue: data[8:12] without prior length check
- Potential bounds issue: data[12:16] without prior length check
- Potential bounds issue: data[16:24] without prior length check
- Potential bounds issue: data[24:32] without prior length check
- Potential bounds issue: data[32:40] without prior length check
- Potential bounds issue: data[40:48] without prior length check
- Missing minimum length validation in DecodeFromBytes

### Field Validation

- Missing version field validation

### Rfc Compliance

- Magic number 5908 should be a named constant
- Magic number 5905 should be a named constant
- Magic number 12 should be a named constant
- Magic number 31 should be a named constant
- Magic number 4330 should be a named constant
- Magic number 1769 should be a named constant
- Magic number 5906 should be a named constant
- Magic number 64 should be a named constant
- Magic number 16 should be a named constant
- Magic number 2016 should be a named constant
- Magic number 1992 should be a named constant
- Magic number 255 should be a named constant
- Magic number 26 should be a named constant
- Magic number 5907 should be a named constant
- Magic number 40 should be a named constant
- Magic number 32 should be a named constant
- Magic number 1305 should be a named constant
- Magic number 1989 should be a named constant
- Magic number 1985 should be a named constant
- Magic number 1361 should be a named constant
- Magic number 2030 should be a named constant
- Magic number 38 should be a named constant
- Magic number 2010 should be a named constant
- Magic number 1119 should be a named constant
- Magic number 128 should be a named constant
- Magic number 24 should be a named constant
- Magic number 958 should be a named constant
- Magic number 791 should be a named constant

### Security

- Potential integer overflow: d * NTP
- Potential integer overflow: d * NTP
- Potential integer overflow: d * NTP
- Potential integer overflow: d * NTP
- Potential integer overflow: d * NTP
- Potential integer overflow: d * NTP

## gre.go

### Bounds Checking

- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1] without prior length check
- Potential bounds issue: data[1] without prior length check
- Potential bounds issue: data[2:4] without prior length check
- Potential bounds issue: buf[0] without prior length check
- Potential bounds issue: buf[1] without prior length check
- Potential bounds issue: buf[0] without prior length check
- Potential bounds issue: buf[0] without prior length check
- Potential bounds issue: buf[0] without prior length check
- Potential bounds issue: buf[0] without prior length check
- Potential bounds issue: buf[0] without prior length check
- Potential bounds issue: buf[1] without prior length check
- Potential bounds issue: buf[0] without prior length check
- Potential bounds issue: buf[1] without prior length check
- Potential bounds issue: buf[1] without prior length check
- Potential bounds issue: buf[2:4] without prior length check
- Potential bounds issue: buf[4:6] without prior length check
- DecodeFromBytes function doesn't check data length
- Missing minimum length validation in DecodeFromBytes

### Field Validation

- Missing version field validation

### Rfc Compliance

- Magic number 32 should be a named constant
- Magic number 2012 should be a named constant
- Magic number 16 should be a named constant
- Magic number 20 should be a named constant
- Magic number 80 should be a named constant
- Magic number 40 should be a named constant

### Security

- Potential integer overflow: uint32 * GRERouting
- Potential integer overflow: Next * GRERouting
- Potential integer overflow: g * GRE
- Potential integer overflow: g * GRE
- Potential integer overflow: g * GRE
- Potential integer overflow: g * GRE
- Potential integer overflow: g * GRE
- Potential infinite loop detected

## vrrp.go

### Bounds Checking

- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1] without prior length check
- Potential bounds issue: data[2] without prior length check
- Potential bounds issue: data[3] without prior length check
- Potential bounds issue: data[4] without prior length check
- Potential bounds issue: data[5] without prior length check
- Potential bounds issue: data[6:8] without prior length check
- Missing minimum length validation in DecodeFromBytes

### Field Validation

- Missing version field validation
- Reserved fields not validated to be zero

### Rfc Compliance

- Magic number 3768 should be a named constant
- Magic number 2338 should be a named constant
- Magic number 16 should be a named constant
- Magic number 2016 should be a named constant
- Magic number 100 should be a named constant

### Security

- Potential integer overflow: v * VRRPv2
- Potential integer overflow: v * VRRPv2
- Potential integer overflow: and * should
- Potential integer overflow: v * VRRPv2
- Potential integer overflow: v * VRRPv2
- Potential integer overflow: v * VRRPv2
- Potential infinite loop detected

## vxlan.go

### Bounds Checking

- Potential bounds issue: buf[1:] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1] without prior length check
- Potential bounds issue: data[1] without prior length check
- Potential bounds issue: data[2:4] without prior length check
- Potential bounds issue: bytes[0] without prior length check
- Potential bounds issue: bytes[1] without prior length check
- Potential bounds issue: bytes[0] without prior length check
- Potential bounds issue: bytes[0] without prior length check
- Potential bounds issue: bytes[1] without prior length check
- Potential bounds issue: bytes[1] without prior length check
- Potential bounds issue: bytes[2:4] without prior length check
- Potential bounds issue: bytes[4:8] without prior length check

### Field Validation

- Reserved fields not validated to be zero

### Rfc Compliance

- Magic number 32 should be a named constant
- Magic number 16 should be a named constant
- Magic number 2016 should be a named constant
- Magic number 7348 should be a named constant
- Magic number 40 should be a named constant
- Magic number 24 should be a named constant
- Magic number 80 should be a named constant

### Security

- Potential integer overflow: vx * VXLAN
- Potential integer overflow: vx * VXLAN
- Potential integer overflow: vx * VXLAN
- Potential integer overflow: vx * VXLAN
- Potential integer overflow: vx * VXLAN

## tls.go

### Bounds Checking

- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1:3] without prior length check
- Potential bounds issue: data[3:5] without prior length check
- DecodeFromBytes function doesn't check data length
- Missing minimum length validation in DecodeFromBytes

### Field Validation

- Missing version field validation

### Rfc Compliance

- Magic number 0304 should be a named constant
- Magic number 0200 should be a named constant
- Magic number 0301 should be a named constant
- Magic number 16 should be a named constant
- Magic number 0300 should be a named constant
- Magic number 5246 should be a named constant
- Magic number 0303 should be a named constant
- Magic number 0302 should be a named constant
- Magic number 2018 should be a named constant

### Security

- Potential integer overflow: t * TLS
- Potential integer overflow: t * TLS
- Potential integer overflow: t * TLS
- Potential integer overflow: t * TLS
- Potential integer overflow: t * TLS
- Potential integer overflow: t * TLS
- Potential integer overflow: t * TLS

## radius.go

### Bounds Checking

- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1] without prior length check
- Potential bounds issue: data[2:4] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1] without prior length check
- Potential bounds issue: data[2:] without prior length check
- Potential bounds issue: data[4:20] without prior length check
- Missing minimum length validation in DecodeFromBytes

### Field Validation

- Reserved fields not validated to be zero

### Rfc Compliance

- Magic number 2866 should be a named constant
- Magic number 2865 should be a named constant
- Magic number 2869 should be a named constant
- Magic number 2868 should be a named constant
- Magic number 2020 should be a named constant
- Magic number 2867 should be a named constant
- Magic number 17 should be a named constant

### Security

- Potential integer overflow: radius * RADIUS
- Potential integer overflow: radius * RADIUS
- Potential integer overflow: radius * RADIUS
- Potential integer overflow: radius * RADIUS
- Potential integer overflow: radius * RADIUS
- Potential integer overflow: radius * RADIUS
- Potential integer overflow: radius * RADIUS
- Potential infinite loop detected

## bfd.go

### Bounds Checking

- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1:] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1:] without prior length check
- Potential bounds issue: data[1:] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1:] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[4:] without prior length check
- Potential bounds issue: data[4:] without prior length check
- Potential bounds issue: data[4:] without prior length check
- Potential bounds issue: data[4:] without prior length check
- Potential bounds issue: data[4:] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[5:] without prior length check
- Potential bounds issue: data[1:5] without prior length check
- Potential bounds issue: data[5:] without prior length check
- Potential bounds issue: data[1:5] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1] without prior length check
- Potential bounds issue: data[2] without prior length check
- Potential bounds issue: data[3] without prior length check
- Potential bounds issue: data[4:] without prior length check
- Potential bounds issue: data[8:] without prior length check
- Potential bounds issue: data[12:] without prior length check
- Potential bounds issue: data[16:] without prior length check
- Potential bounds issue: data[20:] without prior length check
- Potential bounds issue: auth[0] without prior length check
- Potential bounds issue: auth[1] without prior length check
- Potential bounds issue: auth[2] without prior length check
- Potential bounds issue: auth[3:] without prior length check
- Potential bounds issue: auth[3] without prior length check
- Potential bounds issue: auth[4:] without prior length check
- Potential bounds issue: auth[8:] without prior length check
- Potential bounds issue: auth[3] without prior length check
- Potential bounds issue: auth[4:] without prior length check
- Potential bounds issue: auth[8:] without prior length check
- Missing minimum length validation in DecodeFromBytes

### Field Validation

- Missing version field validation
- Reserved fields not validated to be zero

### Rfc Compliance

- Magic number 32 should be a named constant
- Magic number 12 should be a named constant
- Magic number 2010 should be a named constant
- Magic number 2017 should be a named constant
- Magic number 5880 should be a named constant
- Magic number 16 should be a named constant
- Magic number 20 should be a named constant
- Magic number 5881 should be a named constant

### Security

- Potential integer overflow: h * BFDAuthHeader
- Potential integer overflow: AuthHeader * BFDAuthHeader
- Potential integer overflow: d * BFD
- Potential integer overflow: d * BFD
- Potential integer overflow: d * BFD
- Potential integer overflow: d * BFD
- Potential integer overflow: d * BFD
- Potential integer overflow: d * BFD
- Potential integer overflow: d * BFD

## ppp.go

### Bounds Checking

- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1] without prior length check
- Potential bounds issue: bytes[0] without prior length check
- Potential bounds issue: bytes[0] without prior length check
- Potential bounds issue: bytes[1] without prior length check

### Rfc Compliance

- Magic number 16 should be a named constant
- Magic number 100 should be a named constant
- Magic number 2012 should be a named constant

### Security

- Potential integer overflow: p * PPP
- Potential integer overflow: p * PPP
- Potential integer overflow: p * PPP

## pppoe.go

### Bounds Checking

- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[1] without prior length check
- Potential bounds issue: data[2:4] without prior length check
- Potential bounds issue: data[4:6] without prior length check
- Potential bounds issue: bytes[0] without prior length check
- Potential bounds issue: bytes[1] without prior length check
- Potential bounds issue: bytes[2:] without prior length check
- Potential bounds issue: bytes[4:] without prior length check

### Field Validation

- Missing version field validation

### Rfc Compliance

- Magic number 16 should be a named constant
- Magic number 2516 should be a named constant
- Magic number 2012 should be a named constant

### Security

- Potential integer overflow: p * PPPoE
- Potential integer overflow: p * PPPoE

## mpls.go

### Bounds Checking

- Potential bounds issue: data[0] without prior length check
- Potential bounds issue: data[4:] without prior length check

### Rfc Compliance

- Magic number 63 should be a named constant
- Magic number 46 should be a named constant
- Magic number 12 should be a named constant
- Magic number 67 should be a named constant
- Magic number 49 should be a named constant
- Magic number 69 should be a named constant
- Magic number 64 should be a named constant
- Magic number 68 should be a named constant
- Magic number 45 should be a named constant
- Magic number 2012 should be a named constant
- Magic number 65 should be a named constant
- Magic number 40 should be a named constant
- Magic number 61 should be a named constant
- Magic number 32 should be a named constant
- Magic number 48 should be a named constant
- Magic number 62 should be a named constant
- Magic number 100 should be a named constant
- Magic number 66 should be a named constant
- Magic number 47 should be a named constant
- Magic number 60 should be a named constant

### Security

- Potential integer overflow: m * MPLS
- Potential integer overflow: m * MPLS

