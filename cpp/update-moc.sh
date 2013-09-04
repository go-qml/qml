#!/bin/sh

set -e

moc govalue.h -o moc_govalue.pp
moc idletimer.pp -o moc_idletimer.pp
