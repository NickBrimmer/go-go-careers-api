# go-go-careers-api

## A Career Data API Written In Go

To get this API up and going

- Download the repo
- run `make dev`

Current Endpoints:

- `localhost:5000/health` (status)
- `localhost:5000/occupations` (get all occupations)
- `localhost:5000/occupations/13-2051.00` (get occupatoin by id)
- `localhost:5000/occupations/13-2051.00/similar` (get occupatoin by id)
- `localhost:5000/search?q=manager` (search occupations by title)
