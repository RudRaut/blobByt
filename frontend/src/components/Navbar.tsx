import React from 'react';
import { FiGithub } from 'react-icons/fi';

const GITHUB_URL = 'https://github.com/RudRaut'; // <-- Replace with your actual GitHub URL

const Navbar: React.FC = () => (
  <nav className="w-full flex items-center justify-between py-3 px-4">
    <span className="text-xl font-bold text-[#37454d]">Walrus File Storage</span>
    <a
      href={GITHUB_URL}
      target="_blank"
      rel="noopener noreferrer"
      className="flex items-center gap-2 text-[#37454d] hover:text-[#76a0bd] transition-colors font-medium"
    >
      <FiGithub className="text-2xl" />
      <span className="hidden sm:inline">GitHub</span>
    </a>
  </nav>
);

export default Navbar; 