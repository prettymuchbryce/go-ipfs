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

test_expect_success "create two small files" '
  random 1000 7 > afile
  random 1000 9 > bfile
'

test_expect_success "add file using CIDv1 but don't pin" '
  AHASH_CIDv1=$(ipfs add -q --cid-version=1 --raw-leaves=false --pin=false afile)
'

test_expect_success "add file using CIDv0" '
  AHASH_CIDv0=$(ipfs add -q --cid-version=0 afile)
'

test_expect_success "check hashes" '
  test "$(cid-fmt %v-%c $AHASH_CIDv0)" = "cidv0-protobuf" &&
  test "$(cid-fmt %v-%c $AHASH_CIDv1)" = "cidv1-protobuf" &&
  test "$(cid-fmt -v 0 %s $AHASH_CIDv1)" = "$AHASH_CIDv0"
'

test_expect_success "make sure CIDv1 hash really is in the repo" '
  ipfs refs local | grep -q $AHASH_CIDv1
'

test_expect_success "make sure CIDv0 hash really is in the repo" '
  ipfs refs local | grep -q $AHASH_CIDv0
'

test_expect_success "run gc" '
  ipfs repo gc
'

test_expect_success "make sure the CIDv0 hash is in the repo" '
  ipfs refs local | grep -q $AHASH_CIDv0
'

test_expect_success "make sure we can get CIDv0 added file" '
  ipfs cat $AHASH_CIDv0 > thefile &&
  test_cmp afile thefile
'

test_expect_success "make sure the CIDv1 hash is not in the repo" '
  ! ipfs refs local | grep -q $AHASH_CIDv1
'

test_expect_success "clean up" '
  ipfs pin rm $AHASH_CIDv0 &&
  ipfs repo gc &&
  ! ipfs refs local | grep -q $AHASH_CIDv0
'

#
#
#

test_expect_success "add file using CIDv1 but don't pin" '
  ipfs add -q --cid-version=1 --raw-leaves=false --pin=false afile
'

test_expect_success "check that we can access the file when converted to CIDv0" '
  ipfs cat $AHASH_CIDv0 > thefile &&
  test_cmp afile thefile
'

test_expect_success "clean up" '
  ipfs repo gc
'

test_expect_success "add file using CIDv0 but don't pin" '
  ipfs add -q --cid-version=0 --raw-leaves=false --pin=false afile
'

test_expect_success "check that we can access the file when converted to CIDv1" '
  ipfs cat $AHASH_CIDv1 > thefile &&
  test_cmp afile thefile
'

#
#
#

test_expect_success "set up iptb testbed" '
  iptb init -n 2 -p 0 -f --bootstrap=none
'

test_expect_success "start nodes" '
  iptb start &&
  iptb connect 0 1
'

test_expect_success "add afile using CIDv0 to node 0" '
  iptb run 0 ipfs add -q --cid-version=0 afile
'

test_expect_failure "get afile using CIDv1 via node 1" '
  iptb run 1 ipfs --timeout=2s cat $AHASH_CIDv1 > thefile &&
  test_cmp afile thefile
'

test_expect_success "add bfile using CIDv1 to node 0" '
  BHASH_CIDv1=$(iptb run 0 ipfs add -q --cid-version=1 --raw-leaves=false bfile)
'

test_expect_failure "get bfile using CIDv0 via node 1" '
  BHASH_CIDv0=$(cid-fmt -v 0 %s $BHASH_CIDv1)
  iptb run 1 ipfs --timeout=2s cat $BHASH_CIDv0 > thefile &&
  test_cmp bfile thefile
'

test_expect_success "stop testbed" '
  iptb stop
'

test_done
