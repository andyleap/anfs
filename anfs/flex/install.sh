#!/bin/sh

set -o errexit
set -o pipefail

VENDOR=andyleap.net
DRIVER=anfs

driver_dir=$VENDOR${VENDOR:+"~"}${DRIVER}
if [ ! -d "/flexmnt/$driver_dir" ]; then
  mkdir "/flexmnt/$driver_dir"
fi

cp "/driver" "/flexmnt/$driver_dir/.driver"
mv -f "/flexmnt/$driver_dir/.driver" "/flexmnt/$driver_dir/driver"

cp "/$DRIVER" "/flexmnt/$driver_dir/.$DRIVER"
mv -f "/flexmnt/$driver_dir/.$DRIVER" "/flexmnt/$driver_dir/$DRIVER"

while : ; do
  sleep 3600
done
