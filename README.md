# grpc-sample

sample grpc

## tools

- grpcurl
  - install

  ```bash
    brew install grpcurl
  ```

  - usage

  ```bash
    grpcurl -plaintext localhost:8888 list
    grpcurl -plaintext localhost:8888 list pancake.baker.BakePancakeService
    grpcurl -plaintext localhost:8888 pancake.baker.BakePancakeService/Bake # Request
  ```

  - reference
    - <https://qiita.com/yukina-ge/items/a84693f01f3f0edba482>
