import React, { Component } from 'react';
import Form from './Form';
import axios from 'axios';


class App extends Component {
  state = {
    weather_data: []
  }

  componentDidMount() {
    axios.get('./data.json')
      .then(response => {
        this.setState({
          weather_data: response.data.weather_items
        });
      })
    .catch(function(err) {
      console.log(err);
    })  
  }

  render() {
    return (
      <div className="App">
        <h1 className="title">Weather API</h1>
        <Form data={this.state.weather_data}/>
      </div>
    );
  }
}

export default App;
