#!/usr/bin/env bash
#
# Copyright (c) 2017 Jakub Sztandera
# MIT Licensed; see the LICENSE file in this repository.
#

test_description="CID Version 0/1 Duality"

. lib/test-lib.sh

test_init_ipfs

#
#
#

test_expect_success "create a small file" '
  random 1000 7 > afile
'

test_expect_success "add file using CIDv1 but don't pin" '
  HASH1=$(ipfs add -q --cid-version=1 --raw-leaves=false --pin=false afile)
'

test_expect_success "add file using CIDv0" '
  HASH0=$(ipfs add -q --cid-version=0 afile)
'

## note the only difference between these two hashes are the cid
## version the multihash must be the same and they must both be using
## a protobuf
test_expect_success "check hashes" '
  echo $HASH0 $HASH1 &&
  test "$HASH0" = QmUUmUz49Yj3XsPz397V1cFh4wXW7PhYAr8F1YMgw3Qrpb &&
  test "$HASH1" = zdj7Wba1rR812CeQahumt56sz6Fh9iwpw7g3AXmgLkMkUykzX
'

test_expect_success "make sure CIDv1 hash really is in the repo" '
  ipfs refs local | grep -q $HASH1
'

test_expect_success "make sure CIDv0 hash really is in the repo" '
  ipfs refs local | grep -q $HASH0
'

test_expect_success "run gc" '
  ipfs repo gc
'

test_expect_success "make sure the CIDv0 hash is in the repo" '
  ipfs refs local | grep -q $HASH0
'

test_expect_success "make sure we can get CIDv0 added file" '
  ipfs cat $HASH0 > thefile &&
  test_cmp afile thefile
'

test_expect_success "make sure the CIDv1 hash is not in the repo" '
  ! ipfs refs local | grep -q $HASH1
'

test_expect_success "clean up" '
  ipfs pin rm $HASH0 &&
  ipfs repo gc &&
  ! ipfs refs local | grep -q $HASH0
'

#
#
#

test_expect_success "add file using CIDv1 but don't pin" '
  ipfs add -q --cid-version=1 --raw-leaves=false --pin=false afile
'

test_expect_success "check that we can access the file when converted to CIDv0" '
  ipfs cat $HASH0 > thefile &&
  test_cmp afile thefile
'

test_expect_success "clean up" '
  ipfs repo gc
'

test_expect_success "add file using CIDv0 but don't pin" '
  ipfs add -q --cid-version=0 --raw-leaves=false --pin=false afile
'

test_expect_success "check that we can access the file when converted to CIDv1" '
  ipfs cat $HASH1 > thefile &&
  test_cmp afile thefile
'

test_done
