#!/bin/bash

work_dir=`pwd`
package="etcdhooks"
now=`date +%m%d%Y%H%M`
REPO=/var/www/html/repo
tgtdir="$package-$now"

cd $GOPATH

# might fail... so it goes
mkdir -p $GOPATH/bin

# build the binary
go install github.com/tcotav/$package

# test exit code
if [ $? -ne 0 ]
then
  echo "build  of package $package failed"
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

# now push to a web-accessible location to use with automation
mv $tgtdir.tar.gz* $REPO

# remove symlink
rm $REPO/$package-current.tar.gz*

# create new symlinks
ln -s $REPO/$tgtdir.tar.gz $REPO/$package-current.tar.gz
ln -s $REPO/$tgtdir.tar.gz.sha $REPO/$package-current.tar.gz.sha

# run salt or whatever to pick it up -- report error state back to gitlab
