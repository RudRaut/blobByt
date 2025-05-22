import './App.css'

import { Toaster } from 'react-hot-toast'
import Upload from './components/Upload'
import Download from './components/Download'
import StoredFilesList from './components/StoredFilesList'
import Navbar from './components/Navbar'

function App() {
  return (
    <>
      <Navbar />
      <div className="min-h-screen bg-[#f1ede7] flex flex-col items-center py-8 px-2">
        <div className="w-full max-w-4xl bg-white rounded-lg shadow-lg border border-[#e5e7eb] flex flex-col gap-6 p-6">
          <div className="flex flex-col lg:flex-row gap-6 w-full">
            <div className="flex-1 flex flex-col justify-between min-h-[340px]">
              <Upload />
            </div>
            <div className="flex-1 flex flex-col justify-between min-h-[340px]">
              <Download />
            </div>
          </div>
          <div className="w-full">
            <StoredFilesList />
          </div>
        </div>
        <Toaster />
      </div>
    </>
  )
}

export default App
