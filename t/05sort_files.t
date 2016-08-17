#!tapsig

#################
name "Sort lines from multiple files"

cat > input1 <<EOF
2010-05-01T00:00:00Z file 1 line 1
2010-05-01T00:00:02Z file 1 line 2
2010-05-01T00:00:04Z file 1 line 3
EOF

cat > input2 <<EOF
2010-05-01T00:00:01Z file 2 line 1
2010-05-01T00:00:03Z file 2 line 2
2010-05-01T00:00:05Z file 2 line 3
EOF

stdout_is <<EOF
2010-05-01T00:00:01Z file 2 line 1
2010-05-01T00:00:02Z file 1 line 2
2010-05-01T00:00:03Z file 2 line 2
2010-05-01T00:00:04Z file 1 line 3
EOF

tap go-dategrep --from "2010-05-01T00:00:01Z" --to "2010-05-01T00:00:05Z" --format rfc3339 input2 input1

#################
done_testing
