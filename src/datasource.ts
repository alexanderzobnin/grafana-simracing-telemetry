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
        let { telemetry, graph } = target;
        const telemetryField = telemetry || 'Speed';

        const channel = `ds/${this.uid}/${target.source || 'dirtRally2'}`;
        const addr = parseLiveChannelAddress(channel);
        if (!isValidLiveChannelAddress(addr)) {
          continue;
        }

        // const maxLength = request.maxDataPoints ?? 500;
        // Reduce buffer size to improve performance on large dashboards
        const maxLength = graph ? request.maxDataPoints ?? 500 : 2;
        const buffer: StreamingFrameOptions = {
          maxDelta: request.range.to.valueOf() - request.range.from.valueOf(),
          maxLength,
        };

        let filter: any = {
          fields: ['time', telemetryField],
        };
        if (telemetry === '*') {
          // for debugging purposes
          filter = null;
        }

        queries.push(
          getGrafanaLiveSrv().getDataStream({
            key: `${request.requestId}.${counter++}`,
            addr: addr!,
            filter,
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
