import React from 'react';

class Form2 extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            value: '',
            isSelected: false
        };
        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleChange(event) {
        this.setState({
            value: event.target.value
        })
    }

    handleSubmit(event) {
        this.setState({
            isSelected: true
        });
    }

    renderDrawer() {
        if (this.state.isSelected) {
        return (
          <nav>
            <ul>
              <li><a href='#'>Some link</a></li>
              <li><a href='#'>Another link</a></li>
            </ul>
          </nav>
        );
        }
      }

    render() {
        return (
          <div>
            <form onSubmit={this.handleSubmit}>
              <select value="123" onChange={this.handleChange}>
                { this.props.data.map((weather, index) => {
                  return  <option key={weather.woeid} value={weather.woeid}>{weather.title} </option>;
                })}
              </select>
              <input type="submit" value="Submit" />
            </form>
            {this.renderDrawer()}
          </div>
        );        
    }


}

export default Form2;
