interface Sensor {
    id: number;
    Text: string;
    Min: string
    Value: string
    Max: string
    ImageUrl: string
    SensorId?: string,
    Type?: string,
    Children: Sensor[]
}

function makeSensorId(prefix: string, id: string, type: string) {
    const splitId = id.split('/');
    const sensorsWithoutEndingIndex = splitId.slice(0,-1);

    return prefix + '_' + sensorsWithoutEndingIndex + '_' + type.toLowerCase();
}

function makeMetric(sensor: Sensor) {
    return `{host="winnerborn-desktop",objectname="${sensor.Text}"} ${sensor.Value}`
}

function prepare(sensor: Sensor, prefix: string) {
    const metrics = [];

    if (sensor.SensorId && sensor.Type) {
        const sensorId = makeSensorId(prefix, sensor.SensorId, sensor.Type);
        const metric = makeMetric(sensor);

        metrics.push(sensorId, metric);
    } else if (sensor.Children && sensor.Children.length > 0) {
        sensor.Children.forEach((children) => {
            prepare(children, prefix);
        })
    }
}


// cpu_usage_Percent_Processor_Time{host="winnerborn-desktop",instance="0",objectname="Processor",source="winnerborn-desktop"} 7.825819169130066
