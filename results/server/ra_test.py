from result_accumulator import ResultAccumulator

TAG_FILE_PATH = "./results/server/good_tags.yaml"
SCHEMA_FILE_PATH = "./results/server/schema/tags_schema.json"

ra = ResultAccumulator(TAG_FILE_PATH, SCHEMA_FILE_PATH)

b, err = ra.submit_tag("PV003", "Hello")
if err != None:
    print("Tag does not exist:", err)
else:
    print(b)