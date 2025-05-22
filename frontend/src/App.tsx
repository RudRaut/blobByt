import './App.css'

import { Toaster } from 'react-hot-toast'
import Upload from './components/Upload'
import Download from './components/Download'
import StoredFilesList from './components/StoredFilesList'

function App() {

  return (
    <>
      <Upload />
      <Toaster />
      <Download />
      <StoredFilesList />
    </>
  )
}

export default App
