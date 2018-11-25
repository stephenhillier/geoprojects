import axios from 'axios'

const instance = axios.create({
  baseURL: 'https://earthworks.islandcivil.com/api/v1'
});

export default instance
