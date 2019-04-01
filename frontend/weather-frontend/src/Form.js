import React from 'react';

class Form extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            weather_name: '',
            min_temp: '',
            max_temp: '',
            city: '',
            isSelected: false
        };
        this.handleChange = this.handleChange.bind(this);
    }

    handleChange(event) {
        console.log(event.target.value)
        let [name, min_temp, max_temp, city] = event.target.value.split(';');
        this.setState({
            weather_name: name,
            min_temp: parseFloat(min_temp).toFixed(1),
            max_temp: parseFloat(max_temp).toFixed(1),
            city: city,
            isSelected: true
        })
    }

    renderDrawer() {
        if (this.state.isSelected) {
            return (
                <div>
                  <p><b>Forecast: </b>{this.state.weather_name}</p>
                  <p><b>Min temp: </b>{this.state.min_temp} °C</p>
                  <p><b>Max temp: </b>{this.state.max_temp} °C</p>
                  <p><b>City: </b>{this.state.city}</p>
                </div>
            );
        }
      }

    render() {
        return (
          <div>
            <form>
              <select value={this.state.city} onChange={this.handleChange}>
                { this.props.data.map((weather, index) => {
                  return  <option key={weather.woeid} 
                    value={weather.weather_state_name + ';' + weather.min_temp + ';' + weather.max_temp + ';' + weather.title}>
                    {weather.title} </option>;
                })}
              </select>
            </form>
            {this.renderDrawer()}
          </div>
        );        
    }
}

export default Form;
