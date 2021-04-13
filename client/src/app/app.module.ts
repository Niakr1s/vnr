import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { ControlComponent } from './control/control.component';
import { SentenceComponent } from './sentence/sentence.component';
import { TranslationComponent } from './translation/translation.component';
import { ProgressComponent } from './progress/progress.component';
import { LanguageComponent } from './language/language.component';
import { SettingsComponent } from './settings/settings.component';

@NgModule({
  declarations: [
    AppComponent,
    ControlComponent,
    SentenceComponent,
    TranslationComponent,
    ProgressComponent,
    LanguageComponent,
    SettingsComponent,
  ],
  imports: [BrowserModule, AppRoutingModule, HttpClientModule],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
