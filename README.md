# Photo Gallery Generator for Storj DCS

## Requirements

- Golang 1.16 or later
- Account on Storj DCS

## How to create a photo gallery

1. Create a new bucket on Storj DCS.
2. Upload your photos to the `pics/original/` folder of the bucket. The photos must be organized is subfolders. Each subfolder represent an album.
3. Clone this Github repository.
4. Open a Terminal inside the root of the repository.
5. Execute this command:

```
go run main.go generate --access <your-access-grant> --bucket <gallery-bucket>
```

Where:
- `<your-access-grant>` is an access grant to a Storj DCS project
- `<gallery-bucket>` is the bucket on Storj DCS where the photo gallery is hosted

The generator will download the uploaded photos and will upload the generated static website for the photo gallery.
