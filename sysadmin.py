#!/usr/bin/python

import os

filename = "~/.ssh/id_rsa"

if not os.path.exists(filename):
	print("ahh")

instance_list = {}
instance_list['1'] = "Thingworx instance"
instance_list['2'] = "DB instance"
instance_list['3'] = "Services instance"
instance_list['4'] = "Exit"

while True:
	options=instance_list.keys()
	options.sort()
	for entry in options:
		print entry, instance_list[entry]

	selection=raw_input("Please select: ")
	if selection == '1':
		print "add"
	elif selection == '2':
		print "d"
	elif selection == '3':
		print "ass"
	elif selection == '4':
		break
	else:
		print "unknown"
