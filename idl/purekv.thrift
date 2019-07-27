namespace go purekv

service PureService {
  string  Get(1: string key),
  string  Set(1: string key, 2: string value),
  string  Delete(1: string key),
}
