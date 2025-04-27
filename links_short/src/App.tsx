import React from 'react';
import {Layout} from "./components/Layout/Layout";
import {AuthPage} from "./pages/AuthPage/AuthPage";
import {ShortLinkPage} from "./pages/ShortLinkPage/ShortLinkPage";
import {ManagePage} from "./pages/ManageLinksPage/ManagePage";

const App: React.FC = () => {

    return (
        <div className="App">
          {/*@ts-ignore*/}
          <Layout>
            {/*<AuthPage />*/}
              <ShortLinkPage />
            {/*  <ManagePage />*/}
          </Layout>
        </div>
    );
}

export default App;
