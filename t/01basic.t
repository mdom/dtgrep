#!tapsig

#################
name "Empty input results in empty output"

cat > input <<EOF
EOF

tap go-dategrep --from "2010-05-01T00:00:00Z" --to "2010-05-01T00:00:01Z" input

#################
name "Match all input lines"

cat > input <<EOF
2010-05-01T00:00:00Z line 1
2010-05-01T00:00:01Z line 2
2010-05-01T00:00:02Z line 3
EOF

stdout_is <<EOF
2010-05-01T00:00:00Z line 1
2010-05-01T00:00:01Z line 2
2010-05-01T00:00:02Z line 3
EOF

tap go-dategrep --from "2010-05-01T00:00:00Z" --to "2010-05-01T00:00:03Z" --format rfc3339 input

#################
name "Match no lines"

tap go-dategrep --from "2010-05-01T00:00:03Z" --format rfc3339 input

#################
name "Output single line in middle of input"

stdout_is <<EOF
2010-05-01T00:00:01Z line 2
EOF

tap go-dategrep --from "2010-05-01T00:00:01Z" --to "2010-05-01T00:00:02Z" --format rfc3339 input

#################
name "Skip dateless lines"

cat > input <<EOF
2010-05-01T00:00:00Z line 1
foo
2010-05-01T00:00:01Z line 2
2010-05-01T00:00:02Z line 3
EOF

rc_is 1

stderr_is <<EOF
Aborting. Found line without date: foo
EOF

tap go-dategrep --from "2010-05-01T00:00:01Z" --to "2010-05-01T00:00:02Z" --format rfc3339 input

#################
name "Skip dateless lines"

stdout_is <<EOF
2010-05-01T00:00:01Z line 2
EOF

tap go-dategrep --skip-dateless --from "2010-05-01T00:00:01Z" --to "2010-05-01T00:00:02Z" --format rfc3339 input

#################
name "Print multine logs"

cat > input <<EOF
2010-05-01T00:00:00Z line 1
foo
2010-05-01T00:00:01Z line 2
bar
2010-05-01T00:00:02Z line 3
EOF

stdout_is <<EOF
2010-05-01T00:00:01Z line 2
bar
EOF

tap go-dategrep --from "2010-05-01T00:00:01Z" --to "2010-05-01T00:00:02Z"  --format rfc3339 --multiline input

#################
name "Error without format"

cat > input <<EOF
2010-05-01T00:00:01Z line 2
EOF

stderr_is <<EOF
Aborting. Found line without date: 2010-05-01T00:00:01Z line 2
EOF

rc_is 1

tap go-dategrep --from "2010-05-01T00:00:01Z" --to "2010-05-01T00:00:02Z" input

#################
name "Getting format from environment"

cat > input <<EOF
2010-05-01T00:00:01Z line 2
EOF

stdout_is <<EOF
2010-05-01T00:00:01Z line 2
EOF

export GO_DATEGREP_FORMAT=rfc3339
tap go-dategrep --from "2010-05-01T00:00:01Z" --to "2010-05-01T00:00:02Z" input
unset -v GO_DATEGREP_FORMAT


#################

done_testing
