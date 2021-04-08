package transactions

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestCleanUp(t *testing.T) {
	bigPayload := []byte(`{"tipoContrato":{"descricao":"Administrativo","sigla":"ADM","cor":"#00FF00","icone":"tipo_contrato_5.png"},"tipoAdmin":{"descricao":"NãoAplica","sigla":"NAP","cor":"#000000"},"status":{"tipo":"VIGENTE","descricao":"Vigente","sigla":"VIG","cor":"#00FF00","icone":"status_contrato_2.png"},"nrContrato":"00212019","nrContratoExt":"00212019","nrContratoSic":"","nrLicitacao":"001/2019","nrProtocolo":"0559365919","nrMapp":{"nrMapp":0,"descricao":"NãoInformado"},"objeto":"PRESTAÇÃODESERVIÇOSDEINFORMÁTICAPARADISPONIBILIZAÇÃODEINFRAESTRUTURADETIEMNUVEM","observacao":"","contratante":{"pessoaJuridica":{"cnpj":"33866288000130","nome_fantasia":"SOP","site":"","consorcio":false},"tipo":2},"contratada":{"cnpj":"03773788000167","nomeFantasia":"ETICE","site":"","consorcio":false},"interveniente":{"pessoaJuridica":{"cnpj":null,"nome_fantasia":null,"site":null,"consorcio":null},"tipo":null},"distritoOperacional":{"nome":null,"sigla":null,"cidadeSede":{"nome":null,"uf":{"uf":null,"nome":null,"capital":null}}},"dataProposta":"2019-06-19","dataAssinatura":"2019-07-12","dataPublicacao":"2021-02-02","dataOrdemServico":"2021-02-02","prazoInicial":365,"diasAdicionado":0,"diasSuspenso":0,"dataFimVigencia":"2022-02-02","dataFimExecucao":"2021-02-02","totalAditivo":0.00,"totalReajuste":0.00,"valorAtual":450000.00,"valorOriginal":450000.00,"valorPI":450000.00,"aliquotaIRRF":{"tipoImposto":{"sigla":"IRRF","descricao":"ImpostodeRendaRetidonaFonte","especificacao":""},"percentualALIrrf":0.0000,"descricao":"NãoInformado"},"calculoReajuste":{"sigla":"CMA","descricao":"CálculoManual","especificacao":"Cálculonãoserárealizadopelosistema.","cor":"#FF0000"},"exigeGarantia":false,"nrOrdemServico":"34234","calculoImposto":{"sigla":"NAOCAL","descricao":"ImpostoNãoCalculado","especificacao":"Cálculonãorealizadopelosistema.SeráatribuídoatodoscontratosencerradosquenãohouvecontroledeimpostosrealizadospeloSIGDER.","cor":"#FF0000"}}`)

	var obj map[string]interface{}

	err := json.Unmarshal(bigPayload, &obj)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	cleanUp(obj)
	_, err = json.Marshal(obj)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}
