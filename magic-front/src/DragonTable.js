import React from "react";

const URL = "http://localhost:8000/dragons";

export class DragonTable extends React.Component {
  constructor() {
    super();
    this.state = {
      data: null
    }
  }

  componentDidMount() {
    fetch(URL)
      .then(response => response.json())
      .then(json => {
        this.setState({data: json})
      })
  }

  render() {
    if (this.state.data !== null) {
      return (
        <table className="dragonTable">
          <tr>
            <th>Name</th>
            <th>Color</th>
            <th>Mana Cost</th>
            <th>Power</th>
            <th>Toughness</th>
            <th>Ability</th>
          </tr>
        {this.state.data.map(dragon => {
          return (
            <tr>
              <td>{dragon.name}</td>
              <td>{dragon.color}</td>
              <td>{dragon.manaCost}</td>
              <td>{dragon.power}</td>
              <td>{dragon.toughness}</td>
              <td>{dragon.ability}</td>
            </tr>
          )
        })}
        </table>
      )
    }
  }
}
