/** @format */

const echartsConf = {
  backgroundColor: '#212121',
  title: {
    text: '',
    subtext: '',
    x: 'center',
    textStyle: {
      color: '#f2f2f2',
    },
  },
  tooltip: {
    trigger: 'item',
    formatter: '{a} <br/>{b} : {c} ({d}%)',
  },
  legend: {
    orient: 'vertical',
    left: 'left',
    data: ['搜索引擎', '直接访问', '推荐', '其他'],
    textStyle: {
      color: '#f2f2f2',
    },
  },
  series: [
    {
      name: '访问来源',
      type: 'pie',
      radius: '55%',
      center: ['50%', '60%'],
      data: [
        { value: 10440, name: '搜索引擎', itemStyle: { color: '#ef4136' } },
        { value: 4770, name: '直接访问' },
        { value: 2430, name: '推荐' },
        { value: 342, name: '其他' },
        { value: 18, name: '社交平台' },
      ],
      itemStyle: {
        emphasis: {
          shadowBlur: 10,
          shadowOffsetX: 0,
          shadowColor: 'rgba(0, 0, 0, 0.5)',
        },
      },
    },
  ],
}

const echartsDemoJsonStr = JSON.stringify(echartsConf, null, 2)

export default ``
