namespace go purekv

struct BatchOpRequest {
    1: map<string, string> SetParams,
    2: set<string> DelParams,
    3: set<string> GetParams,
}

struct BatchOpResponse {
    1: map<string, string> SetResult,
    2: map<string, string> DelResult,
    3: map<string, string> GetResult,
}

service PureService {
  string  Get(1: string key),
  string  Set(1: string key, 2: string value),
  string  Delete(1: string key),
  map<string, string> MultiGet(1: set<string> keys),
  map<string, string> MultiSet(1: map<string, string> kvs),
  map<string, string> MultiDelete(1: set<string> keys),
  BatchOpResponse BatchOp(1: BatchOpRequest req),
}
