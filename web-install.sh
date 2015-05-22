#!/usr/bin/env bash

GHUSERNAME="Grayda"
DRIVERNAME="sphere-installer"
VERSION="v0.0.1"
FILENAME="sphere-installer_0.0.1_armhf.deb"

echo "$URL"
echo "This script will download $FILENAME and install $DRIVERNAME $VERSION by $GHUSERNAME"
echo -n "Requesting permission to install files to the Sphere. Please enter the password for the 'ninja' account (default password is 'temppwd', minus the quotes)"
sudo with-rw bash
echo "Done!"
echo -n "Downloading sphere-orvibo driver.."
eval wget -P /tmp/ "https://github.com/Grayda/sphere-installer/releases/download/v0.0.1/sphere-installer_0.0.1_armhf.deb"
echo "Done!"
echo -n "Installing package.."
dpkg -i "/tmp/sphere-installer_0.0.1_armhf.deb"
echo "Installation complete. Please visit the labs page to find the driver installer page!"
