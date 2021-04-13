import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { NgxSortableModule } from 'ngx-sortable';
import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { ControlComponent } from './control/control.component';
import { LanguageComponent } from './language/language.component';
import { ProgressComponent } from './progress/progress.component';
import { SentenceComponent } from './sentence/sentence.component';
import { SettingsComponent } from './settings/settings.component';
import { TranslationComponent } from './translation/translation.component';

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
  imports: [
    BrowserModule,
    AppRoutingModule,
    HttpClientModule,
    NgxSortableModule,
  ],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
