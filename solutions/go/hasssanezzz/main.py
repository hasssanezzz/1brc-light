import sys

d = {}
with open(sys.argv[1], 'r') as file:
    content = file.read()
    for enrty in content.split(';'):
        result = enrty.split(',')
        if len(result) != 2:
            continue
        gover, temp = result
        
        temp = int(temp)
        
        if gover in d:
            d[gover]['min'] = min(temp, d[gover]['min'])
            d[gover]['max'] = max(temp, d[gover]['max'])
            d[gover]['sum'] += temp
            d[gover]['visits'] += 1
        else:
            d[gover] = {
                'min': temp, 'max': temp, 'sum': temp, 'visits': 1
            }
        
for gover in d:
    print(f"{gover}={d[gover]['min']}/{d[gover]['max']}/{d[gover]['sum'] // d[gover]['visits']}")