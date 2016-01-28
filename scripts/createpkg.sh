#!/bin/bash

work_dir=`pwd`
package="etcdhooks"
now=`date +%m%d%Y%H%M`
REPO=/var/www/html/repo
tgtdir="$package-$now"

basedir=$work_dir/..

# confirm we're in the right directory
if [ ! -d "$basedir/vendor" ]
then
  echo "something is wrong -- bad basedir: $basedir"
  exit 1
fi

cd $basedir

# build the binary
gb build all

# test exit code
if [ $? -ne 0 ]
then
  echo "build  of package $package failed"
  exit 1
fi

# make deploy dir -- might exist
mkdir -p $basedir/deploy/$tgtdir

# copy over the binaries
cp $basedir/bin/etcdhooks $basedir/deploy/$tgtdir/etcdhooks

# test exit code
if [ $? -ne 0 ]
then
  echo "failed to find binaries"
  exit 1
fi

# then get the scripts we want into the dir
cd $work_dir
cp *.sh  $basedir/deploy/$tgtdir

# lastly the configs we want
cd $basedir
cp etc/*.cfg $basedir/deploy/$tgtdir

# now make archive
cd $basedir/deploy/

tar -czvf $tgtdir.tar.gz $tgtdir
rm -Rf $tgtdir

shasum -a256 $tgtdir.tar.gz > $tgtdir.tar.gz.sha

sudo mkdir -p $REPO

# if we can't make repo, then exit
if [ $? -ne 0 ]
then
  echo "Could not create repo: $REPO"
  exit 1
fi

sudo mv $tgtdir.tar.gz* $REPO

# remove symlink
sudo rm $REPO/$package-current.tar.gz*

# create new symlinks
sudo ln -s $REPO/$tgtdir.tar.gz $REPO/$package-current.tar.gz
sudo ln -s $REPO/$tgtdir.tar.gz.sha $REPO/$package-current.tar.gz.sha

# run salt or whatever to pick it up -- report error state back to gitlab
