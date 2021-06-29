import {
  DataQueryRequest,
  DataQueryResponse,
  DataSourceInstanceSettings,
  isValidLiveChannelAddress,
  parseLiveChannelAddress,
  StreamingFrameOptions,
} from '@grafana/data';
import { MyDataSourceOptions, TelemetryQuery } from './types';

import { DataSourceWithBackend, getGrafanaLiveSrv } from '@grafana/runtime';
import { Observable, of, merge } from 'rxjs';

let counter = 100;

export class DataSource extends DataSourceWithBackend<TelemetryQuery, MyDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
    super(instanceSettings);
  }

  query(request: DataQueryRequest<TelemetryQuery>): Observable<DataQueryResponse> {
    const queries: Array<Observable<DataQueryResponse>> = [];
    for (const target of request.targets) {
      if (target.hide) {
        continue;
      }

      if (target.withStreaming === true || true) {
        let { telemetry } = target;
        const telemetryField = telemetry || 'Speed';

        const channel = `ds/${this.uid}/dirt`
        const addr = parseLiveChannelAddress(channel);
        if (!isValidLiveChannelAddress(addr)) {
          continue;
        }
        const buffer: StreamingFrameOptions = {
          maxLength: request.maxDataPoints ?? 500,
          maxDelta: request.range.to.valueOf() - request.range.from.valueOf(),
        };

        queries.push(
            getGrafanaLiveSrv().getDataStream({
              key: `${request.requestId}.${counter++}`,
              addr: addr!,
              filter: {
                fields: ['time', telemetryField],
              },
              buffer,
            })
        );
      }
    }

    // With a single query just return the results
    if (queries.length === 1) {
      return queries[0];
    }
    if (queries.length > 1) {
      return merge(...queries);
    }
    return of(); // nothing
  }
}
