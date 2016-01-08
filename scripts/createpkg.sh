#!/bin/bash

work_dir=`pwd`
package="etcdhooks"
now=`date +%m%d%Y%H%M`
tgtdir="$package-$now"

cd $GOPATH

# might fail... so it goes
mkdir -p $GOPATH/bin

# build the binary
go install github.com/tcotav/$package

# test exit code
if [ $? -ne 0 ]
then
  echo "go install of package $package failed"
  exit 1
fi

# make deploy dir -- might exist
mkdir -p $GOPATH/deploy/$tgtdir

# copy over the binaries
cp $GOPATH/bin/etcdhooks $GOPATH/deploy/$tgtdir

# then get the scripts we want into the dir
cd $work_dir
cp *.sh  $GOPATH/deploy/$tgtdir

# now make archive
cd $GOPATH/deploy/

tar -czvf $tgtdir.tar.gz $tgtdir
rm -Rf $tgtdir

shasum -a256 $tgtdir.tar.gz $tgtdir.tar.gz.sha  
