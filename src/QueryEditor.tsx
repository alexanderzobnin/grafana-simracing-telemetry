import { defaults } from 'lodash';

import React, { PureComponent, SyntheticEvent } from 'react';
import {InlineField, InlineSwitch, Select} from '@grafana/ui';
import {QueryEditorProps, SelectableValue} from '@grafana/data';
import { DataSource } from './datasource';
import { defaultQuery, MyDataSourceOptions, TelemetryQuery } from './types';
import {dirtRallyOptions} from "./dirtRallyOptions";
import {accOptions} from "./accOptions";
import {iRacingOptions} from "./iRacingOptions";

export const sourceOptions = [
  { label: 'DiRT Rally 2.0', value: 'dirtRally2' },
  { label: 'Assetto Corsa Competizione', value: 'acc' },
  { label: 'iRacing', value: 'iRacing' },
];

type Props = QueryEditorProps<DataSource, TelemetryQuery, MyDataSourceOptions>;

export class QueryEditor extends PureComponent<Props> {
  onTelemetryChange = (option: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, telemetry: option.value });
    // executes the query
    onRunQuery();
  };

  onSourceChange = (option: SelectableValue<string>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, source: option.value });
    onRunQuery();
  };

  onWithStreamingChange = (event: SyntheticEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, withStreaming: event.currentTarget.checked });
    // executes the query
    onRunQuery();
  };

  onGraphChange = (event: SyntheticEvent<HTMLInputElement>) => {
    const { onChange, query, onRunQuery } = this.props;
    onChange({ ...query, graph: event.currentTarget.checked });
    // executes the query
    onRunQuery();
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { telemetry, source, withStreaming, graph } = query;

    let options = dirtRallyOptions;
    if (source === 'acc') {
      options = accOptions;
    } else if (source === 'iRacing') {
      options = iRacingOptions;
    }

    return (
      <div className="gf-form">
        <InlineField label="Source">
          <Select
              width={25}
              options={sourceOptions}
              value={source}
              onChange={this.onSourceChange}
              defaultValue={'acc'}
          />
        </InlineField>
        <Select
          width={25}
          options={options}
          value={telemetry}
          onChange={this.onTelemetryChange}
          defaultValue={'Time'}
        />
        <InlineField label="Enable streaming">
          <InlineSwitch
              value={withStreaming || false}
              onChange={this.onWithStreamingChange}
              css=""
          />
        </InlineField>
        <InlineField label="Graph">
          <InlineSwitch
              value={graph}
              onChange={this.onGraphChange}
              css=""
          />
        </InlineField>
      </div>
    );
  }
}
