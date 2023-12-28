from result_accumulator import ResultAccumulator


ra = ResultAccumulator("./good_tags.yaml", "./tags_schema.json")

b, err = ra.submit_tag("PV003", "Hello")
if err != None:
    print("Tag does not exist:", err)
else:
    print(b)