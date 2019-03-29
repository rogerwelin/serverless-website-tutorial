import React from 'react';


let form_value = '';
let show_desc = false;
let apa = ""

function handleChange(event) {
  form_value = event.target.value;
};

function handleSubmit(event) {
  show_desc = true;
  //alert(show_desc);
  //alert(form_value);
}

if (show_desc) {
  apa = <p>apa</p>
}


const Form = (props) => {
  // console.log(props.data)
  return (
    <div>
      <form onSubmit={handleSubmit}>
        <select value={form_value} onChange={handleChange}>
          { props.data.map((weather, index) => {
           return <option key={weather.woeid} value={weather.woeid}>{weather.title} </option>;
          })}
        </select>
        <input type="submit" value="Submit" />
      </form>
    { apa }
    </div>
  )
}

export default Form;
