import './style'
import React from 'react'
import DailyItem from 'src/components/dailyItem'

const DailyList = props => {
  const { list } = props
  return (
    <article className="daily-list">
      {
        list.map(item => <DailyItem key={item.id} data={item} />)
      }
    </article>
  )
}
export default DailyList