# My solution

### stratigies in my head

#### Idea 1
- spin up workers listening to a channel to recv chunks
- split file into chunks and stream these chunks as they are extracted
- workers will parse these chunks and send entries to a special kind of worker/thread which take these entrie and write it to a hashmap


#### Idea 2
- spin up workers listening to a channel to recv chunks
- split file into chunks and stream these chunks as they are extracted
- workers will parse the chunk and collect these data into a hashMap of a struct called result which will hold `{min, max, sum, count}`
- after all worker are done, the main go routine will compine all of these hashMaps coming from the workers