const CostExplorer = require('aws-sdk/clients/costexplorer')
const SNS = require('aws-sdk/clients/sns')
const moment = require('moment')

const costExplorer = new CostExplorer({
  apiVersion: '2017-10-25',
  region: 'us-east-1'
});
const sns = new SNS({})


/**
 * This method send an SNS notification to your email, witch the forecasted cost
 * for the next month period and the current cost
 */
exports.handler = async () => {
  try {
    const date = new Date();
    const ForecastStartDate = moment(date).add(1, 'd').format('YYYY-MM-DD');
    const ForecastEndDate = moment(date).add(1, 'M').format('YYYY-MM-DD');
    const CurrentStartDate = moment().startOf('month').format('YYYY-MM-DD')
    const CurrentEndDate = moment(date).add(1, 'd').format('YYYY-MM-DD');

    const costAndUsageParams = {
      Metrics: ["BLENDED_COST"],
      Granularity: 'DAILY',
      TimePeriod: {
        End: CurrentEndDate,
        Start: CurrentStartDate
      },
    }

    const forecastParams = {
      Metric: "BLENDED_COST",
      Granularity: 'DAILY',
      TimePeriod: {
        End: ForecastEndDate,
        Start: ForecastStartDate
      },
    }

    const currentCostsResult = await costExplorer.getCostAndUsage(costAndUsageParams).promise()
    const currentTotal = currentCostsResult.ResultsByTime
      .map(result => ({
        total: result.Total.BlendedCost.Amount,
        date: result.TimePeriod.Start
      }))
      .reduce((acc, curr) => {
        if (curr.date === CurrentStartDate) return acc
        return Number(curr.total) + acc
      }, 0)

    const result = await costExplorer.getCostForecast(forecastParams).promise()
    const {Total} = result
    const body = "Today's bill forecast: " + Math.floor(Number(Total.Amount)) + " " + Total.Unit +
      ` for the period ${ForecastStartDate} - ${ForecastEndDate}` +
      `\nCurrent bill: ${Math.floor(currentTotal)} ${Total.Unit} for period ${CurrentStartDate} - ${CurrentEndDate}`;
    console.log(body)

    const snsParams = {
      Message: body,
      TopicArn: process.env.TOPIC_ARN,
      Subject: 'AWS Cost Forecast Notification'
    }

    return await sns.publish(snsParams).promise()
  } catch (e) {
    console.log(e.message)
  }
};
