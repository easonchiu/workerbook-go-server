import './style'
import React from 'react'

const Daily = props => {
  const { data } = props
  return (
    <section className="daily-item">
      <h2>{data.nickname} - {data.groupName}</h2>
      <ul>
        {
          data.dailyList.map(i => (
            <li key={i.id}>
              {i.progress}%{' - '}
              {
                i.pname ? i.pname + ' - ' : null
              }
              {i.record}
            </li>
          ))
        }
      </ul>
    </section>
  )
}
export default Daily