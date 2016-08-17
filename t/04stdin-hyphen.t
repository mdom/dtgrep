#!tapsig

#################
name "Read from stdin with hyphen"

stdout_is <<EOF
2010-05-01T00:00:01Z line 2
EOF

tap go-dategrep --from "2010-05-01T00:00:01Z" --to "2010-05-01T00:00:02Z" --format rfc3339 - <<EOF
2010-05-01T00:00:00Z line 1
2010-05-01T00:00:01Z line 2
2010-05-01T00:00:02Z line 3
EOF

#################
done_testing
