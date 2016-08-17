#!tapsig

depends_on gzip
depends_on bzcat

set -- go-dategrep --from "2010-05-01T00:00:01Z" --to "2010-05-01T00:00:02Z" --format rfc3339

#################
name "Uncompress gzip file on the fly"

gzip > input.gz <<EOF
2010-05-01T00:00:00Z line 1
2010-05-01T00:00:01Z line 2
2010-05-01T00:00:02Z line 3
EOF

stdout_is <<EOF
2010-05-01T00:00:01Z line 2
EOF

tap "$@" input.gz

#################

name "Uncompress bzip2 file on the fly"

bzip2 > input.bz2 <<EOF
2010-05-01T00:00:00Z line 1
2010-05-01T00:00:01Z line 2
2010-05-01T00:00:02Z line 3
EOF

stdout_is <<EOF
2010-05-01T00:00:01Z line 2
EOF

tap "$@" input.bz2

#################
name "Uncompressed and compressed files"

cat > input <<EOF
2010-05-01T00:00:00Z line 1
2010-05-01T00:00:01Z line 2
2010-05-01T00:00:02Z line 3
EOF

gzip > input.gz <<EOF
2010-05-01T00:00:00Z line 1
2010-05-01T00:00:01Z line 2
2010-05-01T00:00:02Z line 3
EOF

stdout_is <<EOF
2010-05-01T00:00:01Z line 2
2010-05-01T00:00:01Z line 2
EOF

tap "$@" input.gz input

#################
done_testing
