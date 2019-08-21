#!/usr/bin/env python
import osquery

# registering all tables
import tables

if __name__ == '__main__':
    osquery.start_extension(name='kubernetes_extension', version='1.0.0')